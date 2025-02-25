FROM golang:1.23.5-alpine AS builder

RUN apk add --no-cache curl

# Yandex Mystem
RUN curl -L -o /tmp/mystem.tar.gz "https://download.cdn.yandex.net/mystem/mystem-3.1-linux-64bit.tar.gz" \
    && tar -xzf /tmp/mystem.tar.gz -C /usr/local/bin/ \
    && rm /tmp/mystem.tar.gz \
    && chmod +x /usr/local/bin/mystem

# Checkup Mystem
RUN /usr/local/bin/mystem -v

WORKDIR /app
COPY . .

# Build app
RUN go build -o main ./cmd/lemmatizer/.

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main /main
COPY --from=builder /usr/local/bin/mystem /usr/local/bin/mystem
ENTRYPOINT ["/main"]
EXPOSE 3000