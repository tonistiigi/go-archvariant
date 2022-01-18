# syntax=docker/dockerfile:1.3

FROM golang:1.17-alpine
RUN apk add --no-cache gcc musl-dev
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.43.0
WORKDIR /go/src/github.com/tonistiigi/binfmt
RUN --mount=target=. --mount=target=/root/.cache,type=cache \
  GOARCH=amd64 golangci-lint run && \
  GOARCH=arm64 golangci-lint run