package regutil

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/opencontainers/go-digest"
	"github.com/spf13/viper"
)

func ListTags(repository string, debug bool) ([]string, error) {
	url, ok := viper.Get("docker_registry_url").(string)
	if !ok {
		msg.Failure("getting tag list: PLATFORM_DOCKER_REGISTRY_URL env is required")
	}
	hub := NewRegistryClient(url, debug)

	return hub.Tags(repository)
}

func DeleteImage(repository, tag string, debug bool) error {
	url, ok := viper.Get("docker_registry_url").(string)
	if !ok {
		msg.Failure("getting tag list: PLATFORM_DOCKER_REGISTRY_URL env is required")
	}
	hub := NewRegistryClient(url, debug)

	s, err := DigestV2(repository, tag)
	if err != nil {
		return err
	}

	digest := digest.Digest(s)
	return hub.DeleteManifest(repository, digest)
}

func NewRegistryClient(url string, debug bool) registry.Registry {
	return registry.Registry{
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
}

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
