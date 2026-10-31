FROM golang:1.25-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o main .

FROM alpine:3.22

WORKDIR /www

LABEL org.opencontainers.image.source=https://github.com/krishankumar01/goravel-taskflow
RUN adduser -D -u 10001 app
COPY --from=builder /build/main ./main
COPY --from=builder /build/public ./public
COPY --from=builder /build/resources ./resources
RUN mkdir -p storage/app storage/logs && chown -R app:app storage
USER app

ENV APP_ENV=production APP_HOST=0.0.0.0 APP_PORT=3000
EXPOSE 3000
ENTRYPOINT ["./main"]
