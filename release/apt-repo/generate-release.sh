#!/bin/sh
set -e

do_hash() {
    HASH_NAME=$1
    HASH_CMD=$2
    echo "${HASH_NAME}:"
    for f in $(find -type f); do
        f=$(echo $f | cut -c3-) # remove ./ prefix
        if [ "$f" = "Release" ]; then
            continue
        fi
        echo " $(${HASH_CMD} ${f}  | cut -d" " -f1) $(wc -c $f)"
    done
}

cat << EOF
Origin: Earthly Technologies
Label: Earthly
Suite: stable
Codename: stable
Version: 1.0
Architectures: amd64 arm64 arm7
Components: main
Description: Earthly software repository
Date: $(date -R)
EOF
do_hash "MD5Sum" "md5sum"
do_hash "SHA1" "sha1sum"
do_hash "SHA256" "sha256sum"
