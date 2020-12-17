#!/usr/bin/env bash
#
# This script will install all released versions of earthly under ~/bin/earthly-v<X.Y.Z>
# It is intended for earthly developers who need to test against previous versions of earthly
# (e.g. making sure a new change doesn't break older versions, or testing out bug reports
# against older versions).
set -e

os="linux"
if [ "$(uname)" == "Darwin" ]; then
    os="darwin"
fi
release_name="earth\\(ly\\)\\?-$os-amd64"

curl -s -L "https://api.github.com/repos/earthly/earthly/releases" > "/tmp/releases.1"

if grep -q 'API rate limit' /tmp/releases.1; then
    if [ ! -f "/tmp/releases" ]; then
        echo "you have been rate-limited by github; and no cached file is available"
        exit 1
    fi
    echo "you have been rate-limited by github; using cached file (if available)"
else
    echo "caching releases under /tmp/releases"
    mv /tmp/releases.1 /tmp/releases
fi

releases=$(jq -r '.[] | @base64' < "/tmp/releases")

for row in $releases; do
    version=$(echo "$row" | base64 -d | jq -r '.name')
    pattern="$version/$release_name"
    urls=$(echo "$row" | base64 -d | jq -r '.assets' | jq -r '.[] | .browser_download_url' | grep "$pattern")
    for  url in $urls; do
        earthlybin="earthly"
        if echo "$url" | grep -w "earth" >/dev/null; then
            earthlybin="earth"
        fi
        outfile="$HOME/bin/$earthlybin-$version"

        if [ ! -f "$outfile" ]; then
            echo "Downloading $url"
            wget -L "$url" -O "$outfile" || ( rm -f "$outfile"; exit 1)
            chmod +x "$outfile"
        else
            echo "$url has already been downloaded to $outfile"
        fi
    done
done
