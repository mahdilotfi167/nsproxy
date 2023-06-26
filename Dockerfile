# syntax=docker/dockerfile:1.4
FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS builder

WORKDIR /code

ENV CGO_ENABLED 0
ENV GOPATH /go
ENV GOCACHE /go-build

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/cache \
    --mount=type=cache,target=/go-build \
    go build -o bin/nsproxy

FROM scratch
COPY --from=builder /code/bin/nsproxy /usr/local/bin/nsproxy

EXPOSE 53/udp

CMD ["/usr/local/bin/nsproxy"]
