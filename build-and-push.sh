#!/bin/bash

echo "Building Docker image..."
docker build -t rogierlommers/perftest .

echo "Pushing Docker image to registry..."
docker push rogierlommers/perftest

echo "Done!"
