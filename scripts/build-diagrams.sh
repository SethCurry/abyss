#!/bin/bash


find . -name "*.puml" -exec plantuml -I./docs/theme.puml-theme {} \;
