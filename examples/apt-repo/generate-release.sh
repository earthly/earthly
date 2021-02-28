#!/bin/sh
set -e

PKGS=$(wc -c ${PACKAGES})
PKGS_GZ=$(wc -c ${PACKAGES}.gz)
cat << EOF
Date: $(date -R)
MD5Sum:
 $(md5sum ${PACKAGES}  | cut -d" " -f1) $PKGS
 $(md5sum ${PACKAGES}.gz  | cut -d" " -f1) $PKGS_GZ
SHA1:
 $(sha1sum ${PACKAGES}  | cut -d" " -f1) $PKGS
 $(sha1sum ${PACKAGES}.gz  | cut -d" " -f1) $PKGS_GZ
SHA256:
 $(sha256sum ${PACKAGES} | cut -d" " -f1) $PKGS
 $(sha256sum ${PACKAGES}.gz | cut -d" " -f1) $PKGS_GZ
EOF
