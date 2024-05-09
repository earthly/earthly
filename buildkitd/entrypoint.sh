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

if [ -z "$CACHE_SIZE_PCT" ]; then
    echo "CACHE_SIZE_PCT not set"
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
    echo "detected cgroups v2; buildkit/entrypoint.sh running under pid=$$ with controllers \"$(cat /sys/fs/cgroup/cgroup.controllers)\" in group $(cat /proc/self/cgroup)"
    test "$(cat /sys/fs/cgroup/cgroup.type)" = "domain" || (echo >&2 "WARNING: invalid root cgroup type: $(cat /sys/fs/cgroup/cgroup.type)")
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
            # Tie-break goes to legacy, after testing on WSL/Windows
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
if [ ! -e "/sbin/$IP_TABLES" ]; then
    echo "IP_TABLES is set to $IP_TABLES, but /sbin/$IP_TABLES does not exist"
    exit 1
fi
ln -sf "/sbin/$IP_TABLES" /sbin/iptables

# clear any leftovers (that aren't explicitly cached) in the dind dir
find /tmp/earthly/dind/ -maxdepth 1 -mindepth 1 | grep -v cache_ | xargs -r rm -rf

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

# Length of time (in seconds) to keep cache. Zero is the same as unset to buildkit.
CACHE_DURATION_SETTINGS=
if [ -n "$CACHE_KEEP_DURATION" ] && [ "$CACHE_KEEP_DURATION" -gt 0 ]; then
  CACHE_DURATION_SETTINGS="$(envsubst </etc/buildkitd.cacheduration.template)"
fi
export CACHE_DURATION_SETTINGS

# For clarity; this will be become CACHE_SIZE_MB after everything is calculated.  It is intentionally left unset
# (and not "0") to simplify the logic.
EFFECTIVE_CACHE_SIZE_MB=

if [ "$CACHE_SIZE_MB" -gt "0" ]; then
    EFFECTIVE_CACHE_SIZE_MB="$CACHE_SIZE_MB"
fi

if [ "$CACHE_SIZE_PCT" -gt "0" ]; then
    # %b -> "Total data blocks"
    # %S -> "Fundamental block size"
    # -f $EARTHLY_TMP_DIR -> filesystem where directory resides, usually a volume in docker's root directory
    CALCULATED_CACHE_MB="$(stat -c "%b * %S * ${CACHE_SIZE_PCT} / 100 / 1024 / 1024" -f "$EARTHLY_TMP_DIR" | bc)"
    if [ -z "$EFFECTIVE_CACHE_SIZE_MB" ]; then
        EFFECTIVE_CACHE_SIZE_MB="$CALCULATED_CACHE_MB"
    elif [ "$CALCULATED_CACHE_MB" -lt "$EFFECTIVE_CACHE_SIZE_MB" ]; then
        echo "clamping cache size to $CALCULATED_CACHE_MB MB (${CACHE_SIZE_PCT}% of filesystem)"
        EFFECTIVE_CACHE_SIZE_MB="$CALCULATED_CACHE_MB"
    else
        # In the off-chance they are actual equal, I'm not sure there's much value in calling that out specifically.
        # Even if they are both "30GB", the user likely set "30000", whereas the percentage would likely come out to
        # be something like "30314" (since we're moving from bytes, unlikely to have a consecutive set of zeroes).
        echo "clamping cache size to fixed size of $EFFECTIVE_CACHE_SIZE_MB MB"
    fi
fi

# EFFECTIVE_CACHE_SIZE_MB remains unset if neither percent nor size were specified.  It would be simpler to just process whether it was
# set (or not), but we'll continue setting to "0" in case anyone has become dependent on that behavior.
CACHE_SIZE_MB="${EFFECTIVE_CACHE_SIZE_MB:-0}"

if [ "$CACHE_SIZE_MB" -eq "0" ]; then
    # no config value was set by the user; buildkit would set this to 10% by default:
    #   https://github.com/moby/buildkit/blob/54b8ff2fc8648c86b1b8c35e5cd07517b56ac2d5/cmd/buildkitd/config/gcpolicy_unix.go#L16
    # however, we will be aggresive and set it to min(55%, max(10%, 20GB))
    CACHE_MB_10PCT="$(stat -c "10 * %b * %S / 100 / 1024 / 1024" -f "$EARTHLY_TMP_DIR" | bc)"
    CACHE_MB_55PCT="$(stat -c "55 * %b * %S / 100 / 1024 / 1024" -f "$EARTHLY_TMP_DIR" | bc)"
    CACHE_SIZE_MB="20480" # first start with 20GB
    if [ "$CACHE_MB_10PCT" -gt "$CACHE_SIZE_MB" ]; then
        CACHE_SIZE_MB="$CACHE_MB_10PCT" # increase it to 10% of the disk if bigger
    elif [ "$CACHE_MB_55PCT" -lt "$CACHE_SIZE_MB" ]; then
        CACHE_SIZE_MB="$CACHE_MB_55PCT" # otherwise, prevent it from being bigger than 55% of the disk
    fi
    echo "cache size set automatically to $CACHE_SIZE_MB MB; this can be changed via the cache_size_mb or cache_size_pct config options"
fi

# Calculate the cache for source files to be 10% of the overall cache
SOURCE_FILE_KEEP_BYTES="$(echo "($CACHE_SIZE_MB * 1024 * 1024 * 0.5) / 1" | bc)" # Note /1 division truncates to int
export SOURCE_FILE_KEEP_BYTES

# convert the cache size into bytes
CATCH_ALL_KEEP_BYTES="$(echo "$CACHE_SIZE_MB * 1024 * 1024" | bc)"
export CATCH_ALL_KEEP_BYTES

# finally populate the cache section of the buildkit toml config
CACHE_SETTINGS="$(envsubst </etc/buildkitd.cache.template)"
export CACHE_SETTINGS

# Set up TCP feature flag, and  also profiling (which has TCP as prerequisite)
TCP_TRANSPORT=
PPROF_SETTINGS=
if [ "$BUILDKIT_TCP_TRANSPORT_ENABLED" = "true" ]; then
    TCP_TRANSPORT="$(cat /etc/buildkitd.tcp.template)"
    if [ "$BUILDKIT_PPROF_ENABLED" = "true" ]; then
        PPROF_SETTINGS="$(cat /etc/buildkitd.pprof.template)"
    fi
fi
export TCP_TRANSPORT
export PPROF_SETTINGS

# Set up TLS feature flag
TLS_ENABLED=
if [ "$BUILDKIT_TLS_ENABLED" = "true" ]; then
    TLS_ENABLED="$(cat /etc/buildkitd.tls.template)"
fi
export TLS_ENABLED

envsubst </etc/buildkitd.toml.template >/etc/buildkitd.toml

# Session history is 1h by default unless otherwise specified
if [ -z "$BUILDKIT_SESSION_HISTORY_DURATION" ]; then
  BUILDKIT_SESSION_HISTORY_DURATION="1h"
fi
export BUILDKIT_SESSION_HISTORY_DURATION

# Session timeout will automatically cancel builds that run for too long
# Configured to 1 day by default unless otherwise specified
if [ -z "$BUILDKIT_SESSION_TIMEOUT" ]; then
  BUILDKIT_SESSION_TIMEOUT="24h"
fi
export BUILDKIT_SESSION_TIMEOUT

# Set up OOM
OOM_SCORE_ADJ="${BUILDKIT_OOM_SCORE_ADJ:-0}"
export OOM_SCORE_ADJ
if [ -n "$OOM_EXCLUDED_PIDS" ]; then
  echo "The following PIDs will be ignored by the OOM reaper: $OOM_EXCLUDED_PIDS"
fi

ignored_by_oom() {
  if echo ",$OOM_EXCLUDED_PIDS," | grep -q ",$1,"; then
    echo "true"
  else
    echo "false"
  fi
}

envsubst "\${OOM_SCORE_ADJ} \${BUILDKIT_DEBUG}" </bin/oom-adjust.sh.template >/bin/oom-adjust.sh
chmod +x /bin/oom-adjust.sh

echo "BUILDKIT_ROOT_DIR=$BUILDKIT_ROOT_DIR"
echo "CACHE_SIZE_MB=$CACHE_SIZE_MB"
echo "BUILDKIT_MAX_PARALLELISM=$BUILDKIT_MAX_PARALLELISM"
echo "BUILDKIT_LOCAL_REGISTRY_LISTEN_PORT=$BUILDKIT_LOCAL_REGISTRY_LISTEN_PORT"
echo "EARTHLY_ADDITIONAL_BUILDKIT_CONFIG=$EARTHLY_ADDITIONAL_BUILDKIT_CONFIG"
echo "CNI_MTU=$CNI_MTU"
echo "OOM_SCORE_ADJ=$OOM_SCORE_ADJ"
echo ""
echo "======== CNI config =========="
cat /etc/cni/cni-conf.json
echo "======== End CNI config =========="
echo ""
echo "======== Buildkitd config =========="
cat /etc/buildkitd.toml
echo "======== End buildkitd config =========="
echo ""
echo "======== OOM Adjust script =========="
cat /bin/oom-adjust.sh
echo "======== OOM Adjust script =========="
echo ""
echo "Detected container architecture is $(uname -m)"

"$@" &
execpid=$!

stop_buildkit() {
  echo "Shutdown signal received. Stopping buildkit..."
  for i in $(echo "$OOM_EXCLUDED_PIDS" | sed "s/,/ /g"); do
    echo "killing externally provided pid: $i"
    kill -SIGTERM "$i"
  done
  echo "killing buildkit pid: $execpid"
  kill -SIGTERM "$execpid"
}

trap stop_buildkit TERM QUIT INT

# quit if buildkit dies
set +x
while true
do
    if ! kill -0 "$execpid" >/dev/null 2>&1; then
        wait "$execpid"
        code="$?"
        if [ "$code" != "0" ]; then
            echo "Error: buildkit process has exited with code $code"
        fi
        exit "$code"
    fi

    for PID in $(pgrep -P 1)
    do
        # Sometimes, child processes can be reparented to the init (this script). One
        # common instance is when something is OOM killed, for instance. This enumerates
        # all those PIDs, and kills them to prevent accidential "ghost" loads.
        if [ "$PID" != "$execpid" ] && [ "$(ignored_by_oom "$PID")" = "false" ]; then
            if [ "$OOM_SCORE_ADJ" -ne "0" ]; then
                ! "$BUILDKIT_DEBUG" || echo "$(date) | $PID($(cat /proc/"$PID"/cmdline)) killed with OOM_SCORE_ADJ=$OOM_SCORE_ADJ" >> /var/log/oom_adj
                kill -9 "$PID"
            else 
                ! "$BUILDKIT_DEBUG" || echo "$(date) | $PID($(cat /proc/"$PID"/cmdline)) was not killed because OOM_SCORE_ADJ was default or not set" >> /var/log/oom_adj
            fi
        fi
    done

    sleep 1
done
