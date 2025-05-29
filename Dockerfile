# Builder
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o blaze

# Runner
FROM alpine AS runner

COPY --from=builder /app/blaze /blaze
COPY --from=builder /app/config.yaml /config.yaml

ENTRYPOINT ["/blaze"]
