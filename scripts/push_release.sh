#!/usr/bin/env bash
set -e

if [ "$#" -ne 2 ]; then
    echo "Some parameters [GIT_TAG, GIT_REPO] were not provided"
    exit
fi

GIT_TAG=$1
GIT_REPO=$2

CHANGELOG=./CHANGELOG.md
CHART=./helm-broker.tar.gz

body="$(cat CHANGELOG.md)"

# Overwrite CHANGELOG.md with JSON data for GitHub API
jq -n \
  --arg body "$body" \
  --arg name "${GIT_TAG}" \
  --arg tag_name "${GIT_TAG}" \
  --arg target_commitish "${GIT_TAG}" \
  '{
    body: $body,
    name: $name,
    tag_name: $tag_name,
    target_commitish: $target_commitish,
    draft: false,
    prerelease: false
  }' > CHANGELOG.md

#CREATE_RELEASE_DATA='{
#  "tag_name": "'"${GIT_TAG}"'",
#  "target_commitish": "'"${GIT_TAG}"'",
#  "name": "'"${GIT_TAG}"'",
#  "body": '$(echo $(cat ${CHANGELOG} | tr -d '\\'))',
#  "draft": false,
#  "prerelease": false
#  }'

echo "Create release GIT_TAG for repo: ${GIT_REPO}, branch: ${GIT_TAG}"
RESPONSE=$(curl -H "Authorization: token ${GITHUB_TOKEN}" --data @CHANGELOG.md "https://api.github.com/repos/${GIT_REPO}/releases")

#RESPONSE=$(curl -s --data "${CREATE_RELEASE_DATA}" "https://api.github.com/repos/$GIT_REPO/releases?access_token=${GITHUB_TOKEN}")
#ASSET_UPLOAD_URL=$(echo "$RESPONSE" | jq -r .upload_url | cut -d '{' -f1)
#if [ -z "$ASSET_UPLOAD_URL" ]; then
#    echo ${RESPONSE}
#    exit 1
#fi
echo ${RESPONSE}

#echo "Uploading CHANGELOG to url: $ASSET_UPLOAD_URL?name=${CHANGELOG}"
#curl -s --data-binary @${CHANGELOG} -H "Content-Type: application/octet-stream" -X POST "$ASSET_UPLOAD_URL?name=$(basename ${CHANGELOG})&access_token=${GITHUB_TOKEN}"
#
#echo "Uploading CHART to url: $ASSET_UPLOAD_URL?name=${CHART}"
#curl -s --data-binary @${CHART} -H "Content-Type: application/octet-stream" -X POST "$ASSET_UPLOAD_URL?name=$(basename ${CHART})&access_token=${GITHUB_TOKEN}"
#
