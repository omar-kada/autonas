#!/usr/bin/env bash
set -euo pipefail

echo "🚫 remove all containers managed by autonas"

ENV_NAME="AUTONAS_MANAGED"
ENV_VALUE="true"

docker ps -aq | while read id; do
    echo "Checking container $id for environment variable $ENV_NAME=$ENV_VALUE"
  # Check if container has the desired environment variable
  if docker inspect -f '{{range .Config.Env}}{{println .}}{{end}}' "$id" | grep -q "^$ENV_NAME=$ENV_VALUE$"; then
    echo "Removing container $id"
    docker rm -f "$id"
  fi
done

echo "✅ all containers removed"