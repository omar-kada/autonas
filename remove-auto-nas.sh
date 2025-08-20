#!/usr/bin/env bash
set -euo pipefail

echo "🚫 remove all containers managed by autonas"

docker ps -q --filter "label=com.autonas.managed=true" | xargs -r docker rm -f

echo "✅ all containers removed"