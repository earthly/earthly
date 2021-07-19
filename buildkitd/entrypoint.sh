#!/bin/sh
set -e
echo "starting earthly-buildkit with EARTHLY_GIT_HASH=$EARTHLY_GIT_HASH BUILDKIT_BASE_IMAGE=$BUILDKIT_BASE_IMAGE"

KERNEL="generic"
if uname -a | grep -wiq "microsoft"; then
  KERNEL="WSL"
  echo "WSL Detected!"
fi
# If we ever need additional kernel detections per-distro, add those checks here and set $KERNEL accordingly.

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

if [ -z "$EARTHLY_TMP_DIR" ]; then
    echo "EARTHLY_TMP_DIR not set"
    exit 1
fi

if [ -z "$NETWORK_MODE" ]; then
    echo "NETWORK_MODE not set"
    exit 1
fi

if [ "$EARTHLY_RESET_TMP_DIR" = "true" ]; then
    echo "Resetting dir $EARTHLY_TMP_DIR"
    rm -rf "${EARTHLY_TMP_DIR:?}"/* || true
fi

if ! lsmod | grep -wq "^ip_tables"; then
  echo "Legacy ip_tables is not loaded."

  if [ "$KERNEL" = "WSL" ]; then
    echo "Keeping iptables-legacy for WSL"

  else
    echo "Switching to iptables-nft."
    ln -sf /sbin/iptables-nft /sbin/iptables
  fi

  lsmod
fi

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

if [ -n "$GIT_URL_INSTEAD_OF" ]; then
    # GIT_URL_INSTEAD_OF can support multiple comma-separated values
    for instead_of in $(echo "${GIT_URL_INSTEAD_OF}" | sed "s/,/ /g")
    do
        base="${instead_of%%=*}"
        insteadOf="${instead_of#*=}"
        git config --global url."$base".insteadOf "$insteadOf"
    done
fi

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
