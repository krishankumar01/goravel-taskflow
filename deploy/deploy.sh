#!/bin/sh
set -eu
cd "$(dirname "$0")"

image="${1:?usage: ./deploy.sh <image>}"
export TASKFLOW_IMAGE="$image"

docker compose pull web

out=$(docker compose run --rm web artisan migrate 2>&1) || true
echo "$out"
if ! echo "$out" | grep -q "Migration success"; then
  echo "migrate failed, keeping the running version"
  exit 1
fi

docker compose up -d --wait
sed -i "s#^TASKFLOW_IMAGE=.*#TASKFLOW_IMAGE=$image#" .env
docker image prune -af
