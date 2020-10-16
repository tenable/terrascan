#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

find_files() {
  find . -not \( \
      \( \
        -wholename '*/vendor/*' \
      \) -prune \
    \) -name '*.go'
}

bad_files=$(find_files | xargs gofmt -d -s 2>&1)
if [[ -n "${bad_files}" ]]; then
  echo "${bad_files}" >&2
  echo >&2
  echo "fix gofmt errors by running:" >&2
  echo "./scripts/fix-gofmt.sh" >&2
  exit 1
fi
