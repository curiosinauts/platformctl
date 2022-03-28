FROM curiosinauts/vscode:0.1.0

COPY --chown=coder:coder  console_linux_amd64  /console
COPY --chown=coder:coder  entrypoint.sh        /entrypoint.sh

RUN  chmod +x /entrypoint.sh             && \
     chmod +x /console              

EXPOSE 22
EXPOSE 3000

USER 1000
ENV USER=coder
WORKDIR /home/coder
USER coder

CMD /entrypoint.sh