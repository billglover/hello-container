FROM scratch
COPY hello /
COPY ./cmd/hello/index.html /
ENTRYPOINT ["/hello"]