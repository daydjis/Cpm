# ---------- build stage ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# зависимости
COPY go.mod go.sum ./
RUN go mod download

# исходники
COPY . .

# сборка
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/app

# ---------- runtime stage ----------
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/app /app/app

CMD ["./app"]
