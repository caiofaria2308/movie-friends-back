FROM golang:1.25-alpine as builder

ARG PROJECT_NAME="app"
ARG VERSION="dev"
ARG COMMIT="none"
ENV CGO_ENABLED=0
ENV GOPROXY=https://proxy.golang.org

RUN apk add curl jq

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app_bin .

FROM debian:stable

WORKDIR /app

RUN apt update -y && apt install curl jq -y

COPY --from=builder /app/app_bin /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT ["/bin/sh", "./docker_assets/entrypoint.sh"]

CMD ["/app/app_bin"]