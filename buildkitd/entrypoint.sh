#!/bin/sh
set -e
echo "starting earthly-buildkit with EARTHLY_GIT_HASH=$EARTHLY_GIT_HASH BUILDKIT_BASE_IMAGE=$BUILDKIT_BASE_IMAGE"

if [ "$BUILDKIT_DEBUG" = "true" ]; then
    set -x
fi

if [ -z "$CACHE_SIZE_MB" ]; then
    echo "CACHE_SIZE_MB not set"
    exit 1
fi

if [ -z "$BUILDKIT_DEBUG" ]; then
    echo "BUILDKIT_DEBUG not set"
    exit 1
fi

if [ -z "$BUILDKIT_MAX_PARALLELISM" ]; then
    echo "BUILDKIT_MAX_PARALLELISM not set"
    exit 1
fi

if [ -z "$EARTHLY_TMP_DIR" ]; then
    echo "EARTHLY_TMP_DIR not set"
    exit 1
fi

if [ -z "$NETWORK_MODE" ]; then
    echo "NETWORK_MODE not set"
    exit 1
fi

if [ -z "$EARTHLY_CACHE_VERSION" ]; then
    echo "EARTHLY_CACHE_VERSION not set"
    exit 1
fi

if [ -f "/sys/fs/cgroup/cgroup.controllers" ]; then
    echo "detected cgroups v2"

    mkdir /sys/fs/cgroup/earthly
    mkdir /sys/fs/cgroup/buildkit
    echo $$ > /sys/fs/cgroup/earthly/cgroup.procs

    echo "+pids" > /sys/fs/cgroup/cgroup.subtree_control
    echo "+cpu" > /sys/fs/cgroup/cgroup.subtree_control

    echo "+pids" > /sys/fs/cgroup/buildkit/cgroup.subtree_control
    echo "+cpu" > /sys/fs/cgroup/buildkit/cgroup.subtree_control

    test "$(cat /sys/fs/cgroup/cgroup.type)" = "domain" || (echo "invalid root cgroup type: $(cat /sys/fs/cgroup/cgroup.type)" && exit 1)
    test "$(cat /sys/fs/cgroup/earthly/cgroup.type)" = "domain" || (echo "invalid earthly cgroup type: $(cat /sys/fs/cgroup/earthly/cgroup.type)" && exit 1)
    test "$(cat /sys/fs/cgroup/buildkit/cgroup.type)" = "domain" || (echo "invalid buildkit cgroup type: $(cat /sys/fs/cgroup/buildkit/cgroup.type)" && exit 1)
fi

earthly_cache_version_path="${EARTHLY_TMP_DIR}/internal.earthly.version"
if [ -f "$earthly_cache_version_path" ]; then
    current_cache_version="$(cat "$earthly_cache_version_path")"
else
    current_cache_version="0"
fi
if [ "$EARTHLY_CACHE_VERSION" != "$current_cache_version" ]; then
    EARTHLY_RESET_TMP_DIR="true"
fi

if [ "$EARTHLY_RESET_TMP_DIR" = "true" ]; then
    echo "Resetting dir $EARTHLY_TMP_DIR"
    rm -rf "${EARTHLY_TMP_DIR:?}"/* || true
    mkdir -p "$EARTHLY_TMP_DIR" # required for eine tests
    echo "$EARTHLY_CACHE_VERSION" > "$earthly_cache_version_path"
fi

if [ -z "$IP_TABLES" ]; then
    echo "Autodetecting iptables"

    if lsmod | grep -wq "^ip_tables"; then
        echo "Detected iptables-legacy module"
        IP_TABLES="iptables-legacy"

    elif lsmod | grep -wq "^nf_tables"; then
        echo "Detected iptables-nft module"
        IP_TABLES="iptables-nft"
    else
        echo "Could not find an ip_tables module; falling back to heuristics."

        legacylines=$(iptables-legacy -t nat -S --wait | wc -l)
        legacycode=$?

        nflines=$(iptables-nft -t nat -S --wait | wc -l)
        nfcode=$?

        if [ $legacycode -eq 0 ] && [ $nfcode -ne 0 ]; then
            echo "Detected iptables-legacy by exit code ($legacycode, $nfcode)"
            IP_TABLES="iptables-legacy"

        elif [ $legacycode -ne 0 ] && [ $nfcode -eq 0 ]; then
            echo "Detected iptables-nft by exit code ($legacycode, $nfcode)"
            IP_TABLES="iptables-nft"

        elif [ $legacycode -ne 0 ] && [ $nfcode -ne 0 ]; then
            echo "iptables-legacy and iptables-nft both exited abnormally ($legacycode, $nfcode). Check your settings and then set the IP_TABLES variable correctly to skip autodetection."
            exit 1

        elif [ "$legacylines" -ge "$nflines" ]; then
            # Tiebreak goes to legacy, after testing on WSL/Windows
            echo "Detected iptables-legacy by output length ($legacylines >= $nflines)"
            IP_TABLES="iptables-legacy"

        else
            echo "Detected iptables-nft by output length ($legacylines < $nflines)"
            IP_TABLES="iptables-nft"
        fi
    fi
else
    echo "Manual iptables specified ($IP_TABLES), skipping autodetection."
fi
ln -sf "/sbin/$IP_TABLES" /sbin/iptables

# clear any leftovers in the dind dir
rm -rf "$EARTHLY_TMP_DIR/dind"
mkdir -p "$EARTHLY_TMP_DIR/dind"

# setup git credentials and config
i=0
while true
do
    varname=GIT_CREDENTIALS_"$i"
    eval data=\$$varname
    # shellcheck disable=SC2154
    if [ -n "$data" ]
    then
        echo 'echo $'$varname' | base64 -d' >/usr/bin/git_credentials_"$i"
        chmod +x /usr/bin/git_credentials_"$i"
    else
        break
    fi
    i=$((i+1))
done
echo "$EARTHLY_GIT_CONFIG" | base64 -d >/root/.gitconfig

#Set up CNI
if [ -z "$CNI_MTU" ]; then
  device=$(ip route show | grep default | cut -d' ' -f5 | head -n 1)
  CNI_MTU=$(cat /sys/class/net/"$device"/mtu)
  export CNI_MTU
fi
envsubst </etc/cni/cni-conf.json.template >/etc/cni/cni-conf.json

# Set up buildkit cache.
export BUILDKIT_ROOT_DIR="$EARTHLY_TMP_DIR"/buildkit
mkdir -p "$BUILDKIT_ROOT_DIR"
CACHE_SETTINGS=
if [ "$CACHE_SIZE_MB" -gt "0" ]; then
    CACHE_SETTINGS="$(envsubst </etc/buildkitd.cache.template)"
fi
export CACHE_SETTINGS

# Set up TCP feature flag
TCP_TRANSPORT=
if [ "$BUILDKIT_TCP_TRANSPORT_ENABLED" = "true" ]; then
    TCP_TRANSPORT="$(cat /etc/buildkitd.tcp.template)"
fi
export TCP_TRANSPORT

# Set up TLS feature flag
TLS_ENABLED=
if [ "$BUILDKIT_TLS_ENABLED" = "true" ]; then
    TLS_ENABLED="$(cat /etc/buildkitd.tls.template)"
fi
export TLS_ENABLED

envsubst </etc/buildkitd.toml.template >/etc/buildkitd.toml
echo "BUILDKIT_ROOT_DIR=$BUILDKIT_ROOT_DIR"
echo "CACHE_SIZE_MB=$CACHE_SIZE_MB"
echo "BUILDKIT_MAX_PARALLELISM=$BUILDKIT_MAX_PARALLELISM"
echo "EARTHLY_ADDITIONAL_BUILDKIT_CONFIG=$EARTHLY_ADDITIONAL_BUILDKIT_CONFIG"
echo "CNI_MTU=$CNI_MTU"
echo ""
echo "======== CNI config =========="
cat /etc/cni/cni-conf.json
echo "======== End CNI config =========="
echo ""
echo "======== Buildkitd config =========="
cat /etc/buildkitd.toml
echo "======== End buildkitd config =========="


echo "Detected container architecture is $(uname -m)"

# start shell repeater server
echo starting shellrepeater
shellrepeater &
shellrepeaterpid=$!

"$@" &
execpid=$!

# quit if either buildkit or shellrepeater die
set +x
while true
do
    if ! kill -0 $shellrepeaterpid >/dev/null 2>&1; then
        echo "Error: shellrepeater process has exited"
        exit 1
    fi
    if ! kill -0 $execpid >/dev/null 2>&1; then
        echo "Error: buildkit process has exited"
        exit 1
    fi
    sleep 1
done
