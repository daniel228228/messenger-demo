#!/bin/bash
move_dir=$1

cd ${move_dir}

if [[ $? -ne 0 ]]; then
  exit 1;
fi

for f in *; do
  if [ -d "$f" ]; then
    find "$f" -maxdepth 1 -mindepth 1 -exec mv {} . \;
    rmdir "$f"
  fi
done