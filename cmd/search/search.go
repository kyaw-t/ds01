package search

import (
	"docker-search/cmd"
	"docker-search/pkg/artifactory"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

type SearchOptions struct {
	Host       string
	SearchTerm string
	Registries []string
	Auth       artifactory.AuthOptions
}

type TagProducerOptions struct {
	Search  SearchOptions
	Wg      *sync.WaitGroup
	Matches <-chan string
	Tags    chan<- string
}

type MatchProducerOptions struct {
	Search  SearchOptions
	Wg      *sync.WaitGroup
	Matches chan<- string
}

type PrintResultOptions struct {
	Results    map[string][]string
	SearchTerm string
	JsonOutput bool
}

var logLevel string
var jsonOutput bool

func runE(_cmd *cobra.Command, args []string) error {
	logger := log.New(log.Writer(), "", 0)
	searchTerm := args[0]
	logger.Print("log-level:", logLevel)
	logger.Print("searching for:", searchTerm)
	results, err := routinedSearch(SearchOptions{
		Host:       "http://localhost:3000",
		SearchTerm: searchTerm,
		Registries: []string{"registry1", "registry2", "registry3", "registry4", "registry5", "registry6", "registry7"},
		Auth: artifactory.AuthOptions{
			Scheme:             "Basic",
			EncodedCredentials: "YWRtaW46cGFzc3dvcmQ=",
		}})

	if err != nil {
		return err
	}

	printResults(PrintResultOptions{
		Results:    results,
		SearchTerm: searchTerm,
		JsonOutput: jsonOutput,
	})

	return nil
}

// naive search implementation for baseline benchmarking
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
func naiveSearch(options SearchOptions) (map[string][]string, error) {
	host := options.Host
	searchTerm := options.SearchTerm
	registries := options.Registries
	resultsM := make(map[string][]string)

	client, err := artifactory.NewArtifactoryClient(host)
	if err != nil {
		return nil, err
	}
	client.Authenticate(options.Auth)

	idx := 0
	for _, registry := range registries {
		repos, err := client.ListDockerRepos(registry, nil, nil)
		if err != nil {
			return nil, err
		}

		for _, repo := range repos {
			if strings.Contains(repo, searchTerm) {
				tags, err := client.ListDockerTags(registry, repo, nil, nil)
				if err != nil {
					return nil, err
				}
				resultsM[registry+"/"+repo] = append(resultsM[registry+"/"+repo], tags...)
				idx += 1
			}
		}
	}
	return resultsM, nil
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "",
	Long:  ``,
	RunE:  runE,
	Args:  cobra.ExactArgs(1),
}

func init() {
	cmd.RootCmd.AddCommand(searchCmd)
	logLevel = *searchCmd.Flags().StringP("log-level", "l", "normal", "[debug | silent]")
	// json Output flag
}

func tagProducer(options TagProducerOptions) {

	defer options.Wg.Done()
	client, err := artifactory.NewArtifactoryClient(options.Search.Host)
	if err != nil {
		log.Fatal(err)
	}
	client.Authenticate(options.Search.Auth)

	for match := range options.Matches {
		parts := strings.Split(match, "@") // registry@repo
		registry, repo := parts[0], parts[1]

		tags, err := client.ListDockerTags(registry, repo, nil, nil)
		if err != nil {
			log.Fatal(err)
		}
		for _, tag := range tags {
			options.Tags <- registry + "/" + repo + ":" + tag
		}
	}
}

func matchProducer(options MatchProducerOptions) {

	defer options.Wg.Done()
	registry := options.Search.Registries[0]
	client, err := artifactory.NewArtifactoryClient("http://localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	client.Authenticate(artifactory.AuthOptions{
		Scheme:             "Basic",
		EncodedCredentials: "YWRtaW46cGFzc3dvcmQ=",
	})
	repos, err := client.ListDockerRepos(registry, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, repo := range repos {
		if strings.Contains(repo, registry) {
			options.Matches <- fmt.Sprintf("%s@%s", registry, repo)
		}
	}
}

func routinedSearch(options SearchOptions) (map[string][]string, error) {
	registries := options.Registries

	var repoWg sync.WaitGroup
	var tagWg sync.WaitGroup
	matches := make(chan string)
	tags := make(chan string)
	maxRoutines := 100

	for _, registry := range registries {
		repoWg.Add(1)
		go matchProducer(MatchProducerOptions{
			Search: SearchOptions{
				Host:       options.Host,
				SearchTerm: registry,
				Registries: []string{registry},
				Auth:       options.Auth,
			},
			Wg:      &repoWg,
			Matches: matches,
		})

	}

	for i := 0; i < maxRoutines; i++ {
		tagWg.Add(1)
		go tagProducer(TagProducerOptions{
			Search:  options,
			Wg:      &tagWg,
			Matches: matches,
			Tags:    tags,
		})
	}

	go func() {
		repoWg.Wait()
		close(matches)
	}()

	go func() {
		tagWg.Wait()
		close(tags)
	}()

	var resultsMap = make(map[string][]string)
	for tag := range tags {
		parts := strings.Split(tag, ":") // registry/repo:tag
		resultsMap[parts[0]] = append(resultsMap[parts[0]], parts[1])
	}

	return resultsMap, nil
}

func printResults(options PrintResultOptions) {

}
