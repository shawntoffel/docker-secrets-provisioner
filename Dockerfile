FROM golang:1.13 as build
WORKDIR /go/src/github.com/shawntoffel/docker-secrets-provisioner
COPY Makefile go.mod go.sum ./
RUN make deps
ADD . .
RUN make build-linux
RUN echo "dsp:x:100:101:/" > passwd

FROM scratch
COPY --from=build /go/src/github.com/shawntoffel/docker-secrets-provisioner/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build --chown=100:101 /go/src/github.com/shawntoffel/docker-secrets-provisioner/bin/dsp .
USER nobody
ENTRYPOINT ["./dsp"]