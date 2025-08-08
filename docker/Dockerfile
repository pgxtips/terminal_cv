# build
FROM golang:latest AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app .

# run
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates openssh-keygen
ENV SSH_HOST_KEY_PATH=/app/keys/ssh_host_ed25519

COPY --from=build /out/app /app/app

# ensure host key exists, then start
COPY docker/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

VOLUME ["/app/keys"]
EXPOSE 1337/tcp
ENTRYPOINT ["/app/entrypoint.sh"]
