# ---- Build ----
FROM golang:1.24.5-bookworm AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/main ./cmd/main.go

RUN apt-get update && apt-get install -y zip && \
    zip -r /build/source_code.zip . -x ".git/*" -x "main" -x "source_code.zip" -x ".env"

# ---- Final Stage ----
FROM gcr.io/distroless/static-debian12

WORKDIR /root/

COPY --from=builder /build/main .
COPY --from=builder /build/source_code.zip .
COPY --from=builder /build/templates ./templates

EXPOSE 8000
CMD ["./main"]