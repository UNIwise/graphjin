# syntax=docker/dockerfile:1

# Base images come through the UNIwise ECR pull-through proxy, like the other service
# images. Outside CI, build with: --build-arg REGISTRY=docker.io/library
ARG REGISTRY=097788601444.dkr.ecr.eu-west-1.amazonaws.com/proxy/library

# stage 1: build. Runs on the build platform and cross-compiles for the target, so a
# multi-arch build does not go through emulation. The web UI is committed as a build
# artifact in serv/web/build, so no Node toolchain is needed here.
FROM --platform=$BUILDPLATFORM ${REGISTRY}/golang:1.27.3-alpine3.24 AS build

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown

WORKDIR /src
COPY . .

# go.work ties the modules together; the binary is the cmd module. timetzdata embeds
# the time zone database so the runtime image needs no tzdata package.
WORKDIR /src/cmd
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -tags timetzdata \
      -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE} -X github.com/dosco/graphjin/serv/v3.version=${VERSION}" \
      -o /out/graphjin .

# stage 2: runtime
FROM ${REGISTRY}/alpine:3.24

# apk upgrade picks up security fixes the base image tag has not been rebuilt with yet.
# The base image already carries the CA bundle.
RUN apk upgrade --no-cache \
 && mkdir -p /config

WORKDIR /
COPY --from=build /out/graphjin /graphjin

ENV GO_ENV=production

ENTRYPOINT ["/graphjin"]
# serv is an alias of serve; the config directory is /config (prod.yaml is mounted there).
CMD ["serv", "--path", "/config"]
