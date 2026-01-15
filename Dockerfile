FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /go/bin/app \
    ./cmd/...

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /go/bin/app /
USER nonroot:nonroot
CMD ["/app"]