#!/usr/bin/env bash

# command:
#   bash release.sh GITHUB_TOKEN

GITHUB_TOKEN=$1

last_tag=$(git describe --tags "$(git rev-list --tags --max-count=1)")

printf "last_tag:%s\n" "${last_tag}"

CREATE_RESPONSE=$(curl \
  -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  https://api.github.com/repos/peng49/gocron/releases \
  -d "{\"tag_name\":\"${last_tag}\",\"target_commitish\":\"master\",\"name\":\"${last_tag}\",\"body\":\"描述\",\"draft\":false,\"prerelease\":false,\"generate_release_notes\":false}")

echo "$CREATE_RESPONSE"

# python3
RELEASE_ID=$(echo "$CREATE_RESPONSE" | python -c "import sys, json; print(json.load(sys.stdin)['id'])")

# upload files
upload_assets() {
  ID=$1
  FILENAME=$2

  printf '\n upload %s\n' "${FILENAME}"

  GH_ASSET="https://uploads.github.com/repos/peng49/gocron/releases/${ID}/assets?name=$(basename "$FILENAME")"

  curl --data-binary @"${FILENAME}" \
    -H "Authorization: token ${GITHUB_TOKEN}" \
    -H "Content-Type: application/octet-stream" "${GH_ASSET}"
}

# shellcheck disable=SC2045
for f in $(ls gocron-package)
do
  upload_assets "${RELEASE_ID}" "gocron-package/${f}"
done

# shellcheck disable=SC2045
for f in $(ls gocron-node-package)
do
  upload_assets "${RELEASE_ID}" "gocron-node-package/${f}"
done