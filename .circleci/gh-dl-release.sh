#!/bin/bash

REPO="vladaionescu/earthly"
VERSION="$1"                       # tag name or the word "latest"
FILE="$2"      # the name of your release asset file, e.g. build.tar.gz
GITHUB="https://api.github.com"

alias errcho='>&2 echo'

function gh_curl() {
  curl -H "Authorization: token $GITHUB_TOKEN" \
       -H "Accept: application/vnd.github.v3.raw" \
       $@
}

if [ "$VERSION" = "latest" ]; then
  # Github should return the latest release first.
  parser=".[0].assets | map(select(.name == \"$FILE\"))[0].id"
else
  parser=". | map(select(.tag_name == \"$VERSION\"))[0].assets | map(select(.name == \"$FILE\"))[0].id"
fi;

asset_id=`gh_curl -s "$GITHUB/repos/$REPO/releases" | jq "$parser"`
if [ "$asset_id" = "null" ]; then
  errcho "ERROR: version not found $VERSION"
  exit 1
fi;

wget -q --auth-no-challenge --header='Accept:application/octet-stream' \
  https://"$GITHUB_TOKEN:@api.github.com/repos/$REPO/releases/assets/$asset_id" \
  -O "$2"
