FROM --platform=$BUILDPLATFORM cgr.dev/chainguard/go:latest-dev AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64

WORKDIR /work

# Copy and download dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the binary with static linking
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o mermaid-check ./cmd/mermaid-check

# Use the static Chainguard image for a minimal, secure, and distroless container
FROM cgr.dev/chainguard/static:latest

# Set labels for better maintainability
LABEL org.opencontainers.image.source="https://github.com/sammcj/mermaid-check"
LABEL org.opencontainers.image.description="A tool to check mermaid diagrams for common issues"
LABEL org.opencontainers.image.licenses="Apache-2.0"

COPY --from=builder /work/mermaid-check /usr/bin/mermaid-check

# Use a non-root user (already provided by Chainguard images, but good for clarity)
USER 65532:65532

ENTRYPOINT ["/usr/bin/mermaid-check"]
