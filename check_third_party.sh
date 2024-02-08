#!/usr/bin/env bash

cd "$(dirname $(realpath $0))"

bash ./web/client/resources/js/lib/update.sh

if ! (git diff-files --quiet web/client/resources/js/lib); then
  echo "Third party libraries don't match their .url files."
  echo "Re-run ./web/client/resources/js/lib/update.sh"
  exit 1
fi
