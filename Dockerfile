# Multi-arch image. Built in CI (Cloud Build) — never locally. Cloud Build
# runners build native amd64 and arm64 in parallel; QEMU is banned.
#
# Image: ghcr.io/hanzoai/template:<tag>
FROM golang:1.26.5-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" \
    -o /out/appd ./cmd/appd

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/appd /usr/local/bin/appd
USER nonroot:nonroot
EXPOSE 8080 9100
ENTRYPOINT ["/usr/local/bin/appd"]
