ARG GOLANG_VERSION=1.23.3

ARG TARGETOS
ARG TARGETARCH

ARG COMMIT
ARG VERSION

FROM --platform=${TARGETARCH} docker.io/golang:${GOLANG_VERSION} AS build

WORKDIR /exporter

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go main.go

COPY collector collector
COPY types types

ARG TARGETOS
ARG TARGETARCH

ARG VERSION
ARG COMMIT

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
    -a -installsuffix cgo \
    -o /go/bin/exporter \
    ./main.go

FROM --platform=${TARGETARCH} gcr.io/distroless/static-debian12:latest

LABEL org.opencontainers.image.description="Prometheus Exporter for Koyeb"
LABEL org.opencontainers.image.source=https://github.com/DazWilkin/koyeb-exporter

COPY --from=build /go/bin/exporter /

EXPOSE 8080

ENTRYPOINT ["/exporter"]
