package regutil

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/docker/distribution/manifest/schema2"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/opencontainers/go-digest"
	"github.com/spf13/viper"
)

// RegistryClient docker registry client
type RegistryClient struct {
	api *registry.Registry
}

// NewRegistryClient returns new docker registry client
func NewRegistryClient(debug bool) (*RegistryClient, error) {
	url, ok := viper.Get("docker_registry_url").(string)
	if !ok {
		return nil, errors.New("error while getting PLATFORMCTL_DOCKER_REGISTRY_URL env value")
	}

	registry := &registry.Registry{
		Logf: func(format string, args ...interface{}) {
			if debug {
				log.Printf(format, args...)
			}
		},
		URL: url,
		Client: &http.Client{
			Transport: http.DefaultTransport,
		},
	}

	return &RegistryClient{api: registry}, nil
}

// ListTags lists docker tags
func (r *RegistryClient) ListTags(repository string, debug bool) ([]string, error) {
	return r.api.Tags(repository)
}

// DeleteImage deletes docker image
func (r *RegistryClient) DeleteImage(repository, tag string, debug bool) error {
	s, err := DigestV2(repository, tag)
	if err != nil {
		return err
	}

	digest := digest.Digest(s)
	return r.api.DeleteManifest(repository, digest)
}

// Repositories lists repositories
func (r *RegistryClient) Repositories() ([]string, error) {
	return r.api.Repositories()
}

// DigestV2 is needed to handling missing functionality from registry.Registry. This client library repo is inactive and
// fails to return correct digest v2.
func DigestV2(repository, reference string) (string, error) {
	base, ok := viper.Get("docker_registry_url").(string)
	if !ok {
		return "", errors.New("DOCKER_REGISTRY_URL missing")
	}

	url := fmt.Sprintf("%s/v2/%s/manifests/%s", base, repository, reference)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", schema2.MediaTypeManifest)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return resp.Header.Get("Docker-Content-Digest"), nil
}
