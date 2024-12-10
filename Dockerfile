ARG GOLANG_VERSION=1.23.3

ARG GOOS=linux
ARG GOARCH=amd64

ARG COMMIT
ARG VERSION

FROM docker.io/golang:${GOLANG_VERSION} as build

WORKDIR /exporter

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go main.go

COPY collector collector
COPY types types

ARG GOOS
ARG GOARCH

ARG VERSION
ARG COMMIT

RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} \
    go build \
    -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
    -a -installsuffix cgo \
    -o /go/bin/exporter \
    ./main.go

FROM gcr.io/distroless/static-debian11:latest

LABEL org.opencontainers.image.description "Prometheus Exporter for Koyeb"
LABEL org.opencontainers.image.source https://github.com/DazWilkin/koyeb-exporter

COPY --from=build /go/bin/exporter /

EXPOSE 8080

ENTRYPOINT ["/exporter"]
