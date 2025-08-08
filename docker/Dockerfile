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
RUN printf '%s\n' \
'#!/bin/sh' \
'set -e' \
'if [ ! -f "$SSH_HOST_KEY_PATH" ]; then' \
'  mkdir -p "$(dirname "$SSH_HOST_KEY_PATH")"' \
'  ssh-keygen -t ed25519 -f "$SSH_HOST_KEY_PATH" -N ""' \
'  echo "Generated host key at $SSH_HOST_KEY_PATH"' \
'fi' \
'echo "Host key fingerprint:"' \
'ssh-keygen -lf "$SSH_HOST_KEY_PATH" || true' \
'exec /app/app -host 0.0.0.0 -port 1337' \
> /app/entrypoint.sh && chmod +x /app/entrypoint.sh

VOLUME ["/app/keys"]
EXPOSE 1337/tcp
ENTRYPOINT ["/app/entrypoint.sh"]
