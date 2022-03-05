ARG LANGUAGES_VERSION
FROM docker-registry.int.curiosityworks.org/7onetella/languages:${LANGUAGES_VERSION}

ENV username=coder

RUN apt-get update --fix-missing  &&    \
    apt-get -q -y  install              \
    `# misc                      `      \
    jq                                  \
    uuid                                \
    zip                                 \
    `# python db dependencies.   `      \
    openssl                             \
    postgresql                          \
    libpq-dev                           \
  && rm -rf /var/lib/apt/lists/*


# Install tools
RUN sudo pip3 install "pip>=20" 


COPY entrypoint.sh /
RUN chmod +x  /entrypoint.sh                            && \
    chown -R ${username}:${username} /home/${username}  


USER 1000
ENV USER=${username}
WORKDIR /home/${username}


EXPOSE 22
CMD [ "/entrypoint.sh" ]