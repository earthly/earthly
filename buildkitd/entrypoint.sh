#!/bin/sh

set -e

if [ -z "$CACHE_SIZE_MB" ]; then
    echo "CACHE_SIZE_MB not set"
    exit 1
fi

if [ -z "$BUILDKIT_DEBUG" ]; then
    echo "BUILDKIT_DEBUG not set"
    exit 1
fi

if [ -z "$EARTHLY_TMP_DIR" ]; then
    echo "EARTHLY_TMP_DIR not set"
    exit 1
fi

if [ "$EARTHLY_RESET_TMP_DIR" == "true" ]; then
    echo "Resetting dir $EARTHLY_TMP_DIR"
    rm -rf "$EARTHLY_TMP_DIR"/* || true
fi

mkdir -p "$EARTHLY_TMP_DIR/dind"

# setup git credentials and config
i=0
while true
do
    varname=GIT_CREDENTIALS_"$i"
    eval data=\$$varname
    if [ -n "$data" ]
    then
        echo 'echo $'$varname' | base64 -d' > /usr/bin/git_credentials_"$i"
        chmod +x /usr/bin/git_credentials_"$i"
    else
        break
    fi
    i=$((i+1))
done
echo "$EARTHLY_GIT_CONFIG" | base64 -d > /root/.gitconfig

if [ -n "$GIT_URL_INSTEAD_OF" ]; then
    # GIT_URL_INSTEAD_OF can support multiple comma-separated values
    for instead_of in $(echo "${GIT_URL_INSTEAD_OF}" | sed "s/,/ /g")
    do
        base="${instead_of%%=*}"
        insteadOf="${instead_of#*=}"
        git config --global url."$base".insteadOf "$insteadOf"
    done

fi

# Set up buildkit cache.
BUILDKIT_ROOT_DIR="$EARTHLY_TMP_DIR"/buildkit
mkdir -p "$BUILDKIT_ROOT_DIR"
echo "BUILDKIT_ROOT_DIR=$BUILDKIT_ROOT_DIR"
echo "CACHE_SIZE_MB=$CACHE_SIZE_MB"
sed 's^:BUILDKIT_ROOT_DIR:^'"$BUILDKIT_ROOT_DIR"'^g; s/:CACHE_SIZE_MB:/'"$CACHE_SIZE_MB"'/g; s/:BUILDKIT_DEBUG:/'"$BUILDKIT_DEBUG"'/g' \
    /etc/buildkitd.toml.template > /etc/buildkitd.toml

echo "ENABLE_LOOP_DEVICE=$ENABLE_LOOP_DEVICE"
echo "FORCE_LOOP_DEVICE=$FORCE_LOOP_DEVICE"
use_loop_device=false
if [ "$FORCE_LOOP_DEVICE" == "true" ]; then
    use_loop_device=true
else
    if [ "$ENABLE_LOOP_DEVICE" == "true" ]; then
        tmp_dir_fs="$(df -T $BUILDKIT_ROOT_DIR | awk '{print $2}' | tail -1)"
        echo "Buildkit dir $BUILDKIT_ROOT_DIR fs type is $tmp_dir_fs"
        if [ "$tmp_dir_fs" != "ext4" ]; then
            echo "Using a loop device, because fs is not ext4"
            use_loop_device=true
        fi
    fi
fi
echo "use_loop_device=$use_loop_device"
if [ "$use_loop_device" == "true" ]; then
    # Create an ext4 fs in a pre-allocated file. Ext4 will allow
    # us to use overlayfs snapshotter even when running on mac.
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
        # We use quadruple the cache size for the loop device. This uses
        # a sparse file for the allocation, meaning that the space is not
        # actually occupied until the cache grows.
        sparse_loop_device_size_mb=$(( CACHE_SIZE_MB * 4 ))
        echo "sparse_loop_device_size_mb=$sparse_loop_device_size_mb"
        dd if=/dev/zero of="$image_file" bs=1M count=0 seek="$sparse_loop_device_size_mb"
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
