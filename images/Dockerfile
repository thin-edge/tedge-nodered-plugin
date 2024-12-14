FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY bin/tedge-nodered-plugin /usr/bin
ENTRYPOINT [ "/usr/bin/tedge-nodered-plugin" ]
CMD ["nodered-flows", "list"]