#!/bin/bash


#find . -name "*.puml" -exec plantuml -I./docs/theme.puml-theme {} \;

for pumlFile in $(find ./site/content -name "*.puml"); do
    echo "Building $pumlFile"
    plantuml -I./docs/theme.puml-theme "$pumlFile"

    filename=$(basename -- "$pumlFile")
    extension="${filename##*.}"
    png_filename="${filename%.*}.png"

    png_file_path="$(dirname -- "$pumlFile")/$png_filename"

    ./scripts/convert-to-jpg.sh "$png_file_path"

    rm "$png_file_path"
done
