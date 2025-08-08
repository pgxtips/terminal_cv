#!/bin/sh
set -e
echo "SSH_HOST_KEY_PATH=${SSH_HOST_KEY_PATH:-/app/keys/ssh_host_ed25519}"

# destruction for all if the key path is a directory
[ -d "$SSH_HOST_KEY_PATH" ] && rm -rf "$SSH_HOST_KEY_PATH"

# generate only if missing/empty (private or .pub)
if [ ! -s "$SSH_HOST_KEY_PATH" ] || [ ! -s "$SSH_HOST_KEY_PATH.pub" ]; then
  mkdir -p "$(dirname "$SSH_HOST_KEY_PATH")"
  rm -f "$SSH_HOST_KEY_PATH" "$SSH_HOST_KEY_PATH.pub"
  ssh-keygen -q -t ed25519 -f "$SSH_HOST_KEY_PATH" -N ""
fi

echo "Host key fingerprint:"
ssh-keygen -lf "$SSH_HOST_KEY_PATH" || true
exec /app/app -host 0.0.0.0 -port 1337
