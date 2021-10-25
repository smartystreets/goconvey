#!/usr/bin/env bash

cd "$(dirname $(realpath $0))"

config_file=()

for url in *.url; do
  config_file+=(-o"$(basename "$url" .url)" url=\"$(cat "$url")\")
done

printf "%s\n" "${config_file[@]}" | curl --parallel --parallel-immediate --config -
