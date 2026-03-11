#!/bin/bash

set -euo pipefail

IMAGE="rogierlommers/perftest"
PLATFORM="linux/amd64"
TAG="$(date +"%Y-%m-%d-%H%M")"

echo "Building and pushing Docker image ${IMAGE}:${TAG} for ${PLATFORM}..."
docker buildx build --platform "${PLATFORM}" -t "${IMAGE}:${TAG}" --push .

echo "Done! Pushed ${IMAGE}:${TAG}"
