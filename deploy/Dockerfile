# Build JS frontend as first step so that it can be bundled with the Go binary
FROM --platform=$BUILDPLATFORM node:21.6.2-alpine3.19 AS front-builder

WORKDIR /build
COPY front/ /build
RUN \
  --mount=type=cache,target=/build/node_modules \
  npm ci && npm run build

# Now build bundled Go binary with both server and client
FROM --platform=$BUILDPLATFORM golang:1.22.0-alpine3.19 AS builder

WORKDIR /build
COPY --from=front-builder /build/dist internal/web/dist
COPY . /build
RUN \
  --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -o loggo ./cmd/loggo

# By default alpine is used for easier debugging,
# but /loggo can be copied from this image to any other linux-based image, e.g. scratch/ubuntu
FROM alpine
COPY --from=builder /build/loggo /bin/loggo

# Basic default configuration for the server
WORKDIR /opt/loggo
EXPOSE 20081/tcp 20080/tcp
VOLUME /opt/loggo

# By default run the hq, since running the client as a separate container is meaningless anyway
ENTRYPOINT [ "/bin/loggo" ]
CMD [ "hq" ]
