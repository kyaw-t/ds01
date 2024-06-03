package artifactory

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func (a *ArtifactoryClient) constructListDockerReposUrl(registry string, limit *int, last *int) (string, error) {
	baseUrl, err := url.Parse(fmt.Sprintf("%s/api/docker/%s/v2/_catalog", a.host, registry))
	if err != nil {
		return "", err
	}

	query := baseUrl.Query()
	if limit != nil {
		query.Set("n", fmt.Sprintf("%d", *limit))
	}
	if last != nil {
		query.Set("last", fmt.Sprintf("%d", *last))
	}

	baseUrl.RawQuery = query.Encode()
	return baseUrl.String(), nil
}

func (a *ArtifactoryClient) constructListDockerTagsUrl(registry string, image string, limit *int, last *int) (string, error) {
	baseUrl, err := url.Parse(fmt.Sprintf("%s/api/docker/%s/v2/%s/tags/list", a.host, registry, image))
	if err != nil {
		return "", err
	}

	query := baseUrl.Query()
	if limit != nil {
		query.Set("n", fmt.Sprintf("%d", *limit))
	}
	if last != nil {
		query.Set("last", fmt.Sprintf("%d", *last))
	}

	baseUrl.RawQuery = query.Encode()
	return baseUrl.String(), nil
}

func (a *ArtifactoryClient) ListDockerRepos(registry string, limit *int, last *int) ([]string, error) {
	if registry == "" {
		return nil, fmt.Errorf("registry is required")
	}
	requestUrl, err := a.constructListDockerReposUrl(registry, limit, last)
	if err != nil {
		return nil, err
	}

	body, err := a.get(requestUrl)
	if err != nil {
		return nil, err
	}

	var response struct {
		Repositories []string `json:"repositories"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Repositories, nil
}

func (a *ArtifactoryClient) ListDockerTags(registry string, image string, limit *int, last *int) ([]string, error) {
	if registry == "" {
		return nil, fmt.Errorf("registry is required")
	}
	requestUrl, err := a.constructListDockerTagsUrl(registry, image, limit, last)
	if err != nil {
		return nil, err
	}

	body, err := a.get(requestUrl)
	if err != nil {
		return nil, err
	}

	var response struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Tags, nil
}
