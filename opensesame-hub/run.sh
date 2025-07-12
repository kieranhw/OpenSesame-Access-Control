#!/bin/bash

OUTPUT_BINARY="opensesame-hub" 
SOURCE_FILE="cmd/opensesame/main.go"

go build -o $OUTPUT_BINARY $SOURCE_FILE

if [ $? -ne 0 ]; then
  echo "Build failed! Exiting."
  exit 1
fi

./$OUTPUT_BINARY
