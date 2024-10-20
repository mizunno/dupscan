#!/bin/bash

DIR="./files"
NUM_FILES=1000
FILE_SIZE=1024  # Size in bytes

# Create directory if it doesn't exist
mkdir -p "$DIR"

# Create files
for i in $(seq 1 $NUM_FILES); do
  dd if=/dev/urandom of="$DIR/file_$i.dat" bs=$FILE_SIZE count=1 status=none
  echo "Created file_$i.dat"
done

