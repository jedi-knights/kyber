# Stage 1: build the static binary.
FROM golang:1.23-alpine AS builder

WORKDIR /src

# Cache module downloads separately from source changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)" \
    -o /out/kyber \
    ./cmd/kyber

# Stage 2: minimal runtime image.
FROM gcr.io/distroless/static-debian12

COPY --from=builder /out/kyber /usr/local/bin/kyber

ENTRYPOINT ["/usr/local/bin/kyber"]
CMD ["analyze", "./..."]
