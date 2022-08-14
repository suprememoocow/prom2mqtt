# Invoked from goreleaser, uses binaries build by goreleaser
FROM alpine:3.16
ENTRYPOINT ["/usr/local/bin/prom2mqtt"]
COPY prom2mqtt /usr/local/bin
