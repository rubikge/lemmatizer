FROM golang:1.23.5-alpine AS builder

RUN apk add --no-cache curl

# Yandex Mystem
RUN curl -L -o /tmp/mystem.tar.gz "https://download.cdn.yandex.net/mystem/mystem-3.1-linux-64bit.tar.gz" \
    && tar -xzf /tmp/mystem.tar.gz -C /usr/local/bin/ \
    && rm /tmp/mystem.tar.gz \
    && chmod +x /usr/local/bin/mystem

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

# Build app
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates libc6-compat
COPY --from=builder /app/main /main
COPY --from=builder /usr/local/bin/mystem /usr/local/bin/mystem
ENTRYPOINT ["/main"]
EXPOSE 3000