package artifactory

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type AuthOptions struct {
	Scheme             string // supports Bearer and Basic
	Token              string // required if scheme is Bearer
	EncodedCredentials string // required if scheme is Basic
}

type ArtifactoryClient struct {
	host    string
	auth    string
	timeout time.Duration
}

func validateHost(host string) error {
	if host == "" {
		return fmt.Errorf("host is required")
	}
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		return fmt.Errorf("host must start with http:// or https://")
	}
	return nil
}

func EncodeBasicCredentials(username string, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, username)))
}

func NewArtifactoryClient(host string) (*ArtifactoryClient, error) {
	err := validateHost(host)
	if err != nil {
		return nil, err
	}
	host = strings.TrimSuffix(host, "/")

	return &ArtifactoryClient{
		host:    host,
		timeout: 5 * time.Second,
	}, nil
}

func (a *ArtifactoryClient) Authenticate(auth AuthOptions) error {

	var authorization string

	switch auth.Scheme {
	case "Bearer":
		if auth.Token == "" {
			return fmt.Errorf("token is required for Bearer authentication")
		}
		authorization = fmt.Sprintf("Bearer %s", auth.Token)

	case "Basic":
		if auth.EncodedCredentials == "" {
			return fmt.Errorf("credentials required for Basic authentication")
		}
		encoded := auth.EncodedCredentials
		authorization = fmt.Sprintf("Basic %s", encoded)

	case "":
		return fmt.Errorf("authentication scheme is required")

	default:
		return fmt.Errorf("unsupported authentication scheme: %s", auth.Scheme)
	}

	a.auth = authorization
	return nil
}

func (a *ArtifactoryClient) get(url string) ([]byte, error) {

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if a.auth != "" {
		req.Header.Add("Authorization", a.auth)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
