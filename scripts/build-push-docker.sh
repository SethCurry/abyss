#!/bin/bash

usage() { echo "Usage: $0 [-p <string>] [-i <string>]" 1>&2; exit 1; }

do_push=0
tag_name=""
image_filter=""

while getopts ":p:i:" o; do
    case "${o}" in
        p)
            do_push=1
            tag_name="${OPTARG}"
            ;;
        i)
            image_filter="${OPTARG}"
            ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

if [ -n "$image_filter" ] && [ ! -d "./build/docker/$image_filter" ]; then
  echo "No such image: ./build/docker/$image_filter" 1>&2
  exit 1
fi

for i in ./build/docker/*; do
  image_name=$(basename $i)

  if [ -n "$image_filter" ] && [ "$image_name" != "$image_filter" ]; then
    continue
  fi

  if [ "$image_name" == "base-js" ]; then
    latest_url="abyss-$image_name:latest"
    dev_url="abyss-$image_name:dev"
    repo_url="abyss-$image_name:release"
  else
    if [ "$do_push" -eq 1 ]; then
      repo_url="ghcr.io/sethcurry/abyss-$image_name:$tag_name"
      latest_url="ghcr.io/sethcurry/abyss-$image_name:latest"
    else
      repo_url="ghcr.io/sethcurry/abyss-$image_name:dev"
      latest_url="ghcr.io/sethcurry/abyss-$image_name:latest"
      dev_url="ghcr.io/sethcurry/abyss-$image_name:dev"
    fi
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
