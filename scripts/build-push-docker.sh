#!/bin/bash

for i in ./build/docker/*; do
  image_name=$(basename $i)
  repo_url="ghcr.io/sethcurry/abyss-$image_name:latest"

  echo "Building $repo_url"
  docker buildx build -t "$repo_url" -f ./build/docker/$image_name/Dockerfile .

  echo "Pushing $repo_url"
  docker push "$repo_url"
done
