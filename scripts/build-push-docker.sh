#!/bin/bash

usage() { echo "Usage: $0 [-s <45|90>] [-p <string>]" 1>&2; exit 1; }

do_push=0

while getopts ":p:" o; do
    case "${o}" in
        p)
            do_push=1
            ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

for i in ./build/docker/*; do
  image_name=$(basename $i)
  repo_url="ghcr.io/sethcurry/abyss-$image_name:latest"

  echo "Building $repo_url"
  docker buildx build -t "$repo_url" -f ./build/docker/$image_name/Dockerfile .

  if [ "$do_push" -eq 1 ]; then
    echo "Pushing $repo_url"
    docker push "$repo_url"
  fi
done
