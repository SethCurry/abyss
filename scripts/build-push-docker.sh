#!/bin/bash

usage() { echo "Usage: $0 [-p <string>]" 1>&2; exit 1; }

do_push=0
tag_name=""

while getopts ":p:" o; do
    case "${o}" in
        p)
            do_push=1
            tag_name="$2"
            ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

for i in ./build/docker/*; do
  image_name=$(basename $i)
  latest_url=""
  dev_url="ghcr.io/sethcurry/abyss-$image_name:dev"

  if [ "$do_push" -eq 1 ]; then
    repo_url="ghcr.io/sethcurry/abyss-$image_name:$tag_name"
    latest_url="ghcr.io/sethcurry/abyss-$image_name:latest"
  else
    repo_url="abyss-$image_name:dev"
    latest_url="abyss-$image_name:latest"
  fi

  echo "Building $repo_url"
  docker buildx build --pull=false -t "$dev_url" -t "$repo_url" -t "$latest_url" -f ./build/docker/$image_name/Dockerfile .

  if [ "$do_push" -eq 1 ]; then
    echo "Pushing $repo_url"
    docker push "$repo_url"
    echo "Pushing $latest_url"
    docker push "$latest_url"
  fi
done
