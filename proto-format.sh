#!/usr/bin/env bash

for path in *.proto; do
    echo " ==> $path"
    sed -i 's/repeated /repeated_/g' "$path" &&
    # uses clang-format from root directory
    clang-format -i -style=file "$path" &&
    sed -i 's/repeated_/repeated /g' "$path"
done

exit 0

