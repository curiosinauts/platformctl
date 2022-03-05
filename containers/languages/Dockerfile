ARG BASE_VERSION
FROM docker-registry.int.curiosityworks.org/7onetella/base:${BASE_VERSION}


RUN apt-get update --fix-missing  &&    \
    apt-get -q -y  install              \
    python3                             \
    python3-pip                         \
  && rm -rf /var/lib/apt/lists/*


ENV GO_VERSION 1.17.6
RUN curl -s -O https://dl.google.com/go/go"${GO_VERSION}".linux-amd64.tar.gz    && \
    tar -xf go"${GO_VERSION}".linux-amd64.tar.gz -C /usr/local                  && \
    rm go"${GO_VERSION}".linux-amd64.tar.gz                                     