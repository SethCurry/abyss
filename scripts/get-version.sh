#!/bin/bash

tag=$(git tag --points-at HEAD)
if [ -n "$tag" ]; then
  echo "$tag"
else
  date -u +%Y-%m-%dT%H:%M:%SZ
fi
