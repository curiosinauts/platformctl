package regutil

import (
	"log"
	"net/http"

	"github.com/curiosinauts/platformctl/internal/msg"
	"github.com/heroku/docker-registry-client/registry"
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

	manifest, err := hub.ManifestV2(repository, tag)
	digest := manifest.Config.Digest
	if err != nil {
		return err
	}

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
