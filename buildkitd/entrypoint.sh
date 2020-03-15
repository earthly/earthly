#!/bin/sh

set -e

if [ -z "$CACHE_SIZE_MB" ]; then
    echo "CACHE_SIZE_MB not set"
    exit 1
fi

if [ -z "$EARTHLY_TMP_DIR" ]; then
    echo "EARTHLY_TMP_DIR not set"
    exit 1
fi

BUILDKIT_ROOT_DIR="$EARTHLY_TMP_DIR"/buildkit
# Leave 1GB as additional buffer.
buildkit_cache_size_mb=$(( CACHE_SIZE_MB - 1000 ))
sed 's^:BUILDKIT_ROOT_DIR:^'"$BUILDKIT_ROOT_DIR"'^g; s/:CACHE_SIZE_MB:/'"$buildkit_cache_size_mb"'/g' \
    /etc/buildkitd.toml.template > /etc/buildkitd.toml

if [ -n "$GIT_URL_INSTEAD_OF" ]; then
    # TODO: Perhaps allow multiple such values to be passed in (eg comma-separated).
    base="${GIT_URL_INSTEAD_OF%%=*}"
    insteadOf="${GIT_URL_INSTEAD_OF#*=}"
    git config --global url."$base".insteadOf "$insteadOf"
fi

if [ -n "$DOCKER_CONFIG_JSON" ]; then
    if [ -z "$DOCKER_CONFIG" ]; then
        if [ -z "$HOME" ]; then
            HOME=/root
        fi
        DOCKER_CONFIG="$HOME/.docker"
    fi
    printenv DOCKER_CONFIG_JSON > "$DOCKER_CONFIG/config.json"
fi

# Create an ext4 fs in a pre-allocated file. Ext4 will allow
# us to use overlayfs snapshotter even when running on mac.
if [ "$ENABLE_LOOP_DEVICE" == "true" ]; then
    echo "ENABLE_LOOP_DEVICE=true"
    echo "CACHE_SIZE_MB=$CACHE_SIZE_MB"
    image_file="$EARTHLY_TMP_DIR"/buildkit.img
    mount_point="$BUILDKIT_ROOT_DIR"

    function do_mount {
        echo "Mounting loop device"
        ret=0
        mount -n -o loop,noatime,nodiratime,noexec,noauto "$image_file" "$mount_point" || ret=1
        return "$ret"
    }

    function init_mount {
        echo "Creating loop device"
        mkdir -p "$mount_point"
        fallocate -l "$CACHE_SIZE_MB"M "$image_file" || \
            dd if=/dev/zero of="$image_file" bs=1M count=0 seek="$CACHE_SIZE_MB"
        mkfs.ext4 "$image_file"
    }

    function reset_mount {
        echo "Resetting loop device"
        umount "$mount_point" || true
        rm -rf "$image_file"
        rm -rf "$mount_point"
    }

    if [ -f "$image_file" ]; then
        do_mount || (reset_mount && init_mount && do_mount)
    else
        init_mount
        do_mount
    fi
fi

exec "$@"
