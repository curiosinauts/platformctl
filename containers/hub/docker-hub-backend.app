[docker-hub]
port   = 80
memory = 128
health = /

env.REGISTRY_URL    = https://docker-registry.int.curiosityworks.org
env.DELETE_IMAGES   = true
env.SINGLE_REGISTRY = true

tag.1               = urlprefix-docker-hub.7onetella.net