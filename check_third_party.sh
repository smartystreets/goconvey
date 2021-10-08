#!/usr/bin/env bash

./web/client/resources/js/lib/update.sh

if ! (git diff-files --quiet); then
  echo "Third party libraries don't match their .url files."
  echo "Re-run ./web/client/resources/js/lib/update.sh"
  exit 1
fi
