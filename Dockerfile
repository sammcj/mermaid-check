FROM cgr.dev/chainguard/go:latest-dev as builder
WORKDIR /work

COPY go.mod /work/
COPY cmd /work/cmd
COPY internal /work/internal

RUN CGO_ENABLED=0 go build -o mermaid-check ./cmd

# for static Go builds
FROM cgr.dev/chainguard/static
# for dynamic Go builds
#FROM cgr.dev/chainguard/glibc-dynamic
COPY --from=builder /work/mermaid-check /mermaid-check

ENTRYPOINT ["/mermaid-check"]
