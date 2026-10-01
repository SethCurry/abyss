#!/bin/bash

file_to_convert="$1"

filename=$(basename -- "$file_to_convert")
extension="${filename##*.}"
new_filename="${filename%.*}.jpg"

new_file_path="$(dirname -- "$file_to_convert")/$new_filename"

magick -quality 80 "$file_to_convert" "$new_file_path"
