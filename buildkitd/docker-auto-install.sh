#!/bin/sh

set -eu

distro=$(. /etc/os-release && echo "$ID")
DOCKER_VERSION="${DOCKER_VERSION:-}"

detect_dockerd() {
    set +e
    command -v dockerd >/dev/null
    has_d="$?"
    set -e
    return "$has_d"
}

detect_docker_compose() {
    set +e
    command -v docker-compose >/dev/null
    has_dc="$?"
    set -e
    return "$has_dc"
}

detect_docker_compose_cmd() {
    if command -v docker-compose >/dev/null; then
        echo "docker-compose"
        return 0
    fi
    if docker help | grep -w compose >/dev/null; then
        echo "docker compose"
        return 0
    fi
    echo >&2 "failed to detect docker compose / docker-compose command"
    return 1
}

detect_jq() {
    set +e
    command -v jq >/dev/null
    has_jq="$?"
    set -e
    return "$has_jq"
}

print_debug() {
    set +u
    if [ "$EARTHLY_DEBUG" = "true" ] ; then
        echo "$@"
    fi
    set -u
}

detect_alpine_3_18_or_newer() {
    VERSION="$(. /etc/os-release && echo "$VERSION_ID")"
    if [ -z "$VERSION" ]; then
        echo >&2 "Error: unable to detect alpine version"
        exit 1
    fi
    MAJOR="$(echo "$VERSION" | awk -F. '{print $1}')"
    MINOR="$(echo "$VERSION" | awk -F. '{print $2}')"
    if [ "$MAJOR" -lt 3 ]; then
        return 1
    fi
    if [ "$MINOR" -lt 18 ]; then
        return 1
    fi
    return 0
}

install_docker_compose() {
    case "$distro" in
        alpine)
            if detect_alpine_3_18_or_newer; then
                apk add --update --no-cache docker-cli-compose
            else
                apk add --update --no-cache docker-compose
            fi
            ;;
        *)
            echo "Detected architecture is $(uname -m)"
            case "$(uname -m)" in
                armv7l|armhf)
                    curl -L "https://github.com/linuxserver/docker-docker-compose/releases/download/1.27.4-ls27/docker-compose-armhf" -o /usr/local/bin/docker-compose
                    ;;
                arm64|aarch64)
                    curl -L "https://github.com/linuxserver/docker-docker-compose/releases/download/1.27.4-ls27/docker-compose-arm64" -o /usr/local/bin/docker-compose
                    ;;
                *)
                    curl -L "https://github.com/docker/compose/releases/download/1.27.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
                    ;;
            esac
            chmod +x /usr/local/bin/docker-compose
            ;;
    esac
}

install_dockerd() {
    case "$distro" in
        alpine)
            if [ -n "$DOCKER_VERSION" ]; then
              apk add --update --no-cache docker="$DOCKER_VERSION"
            else
              apk add --update --no-cache docker
            fi
            ;;

        amzn)
            install_dockerd_amazon
            ;;

        ubuntu)
            install_dockerd_debian_like
            ;;

        debian)
            install_dockerd_debian_like
            ;;

        *)
            echo "Warning: Distribution $distro not yet supported for Docker-in-Earthly."
            echo "Will attempt to treat like Debian."
            echo "If you would like this distribution to be supported, please open a GitHub issue: https://github.com/earthly/earthly/issues"
            install_dockerd_debian_like
            ;;
    esac
}

apt_update_done="false"
apt_get_update() {
    if [ "$apt_update_done" != "true" ]; then
        apt-get update
        apt_update_done=true
    fi
}

install_docker_apt_repo_old() {
    curl -fsSL "https://download.docker.com/linux/$distro/gpg" | apt-key add -
    add-apt-repository \
        "deb [arch=$(dpkg --print-architecture)] https://download.docker.com/linux/$distro \
        $(lsb_release -cs) \
        stable"
}

install_docker_apt_repo_new() {
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL "https://download.docker.com/linux/$distro/gpg" | gpg --no-tty --dearmor -o /etc/apt/keyrings/docker.gpg
    chmod a+r /etc/apt/keyrings/docker.gpg
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/$distro \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      tee /etc/apt/sources.list.d/docker.list > /dev/null
}

install_dockerd_debian_like() {
    export DEBIAN_FRONTEND=noninteractive
    apt-get remove -y docker docker-engine docker.io containerd runc || true
    apt_get_update
    apt-get install -y \
        apt-transport-https \
        ca-certificates \
        curl \
        gnupg-agent \
        software-properties-common

    VERSION="$(. /etc/os-release && echo "$VERSION_ID")"
    case "$distro" in
        ubuntu)
            MAJOR="$(echo "$VERSION" | awk -F. '{print $1}')"
            if [ "$MAJOR" -ge "22" ]; then
                install_docker_apt_repo_new
            else
                install_docker_apt_repo_old
            fi
            ;;

        debian)
            if [ "$VERSION" -ge "12" ]; then
                install_docker_apt_repo_new
            else
                install_docker_apt_repo_old
            fi
            ;;

        *)
            install_docker_apt_repo_old
            ;;
    esac
    apt-get update # dont use apt_get_update since we must update the newly added apt repo
    if [ -n "$DOCKER_VERSION" ]; then
        apt-get install -y docker-ce="$DOCKER_VERSION" docker-ce-cli="$DOCKER_VERSION"
    else
        apt-get install -y docker-ce docker-ce-cli
    fi
    apt-get install -y containerd.io
}

install_dockerd_amazon() {
    version=$(. /etc/os-release && echo "$VERSION")
    case "$version" in
        2)
            yes | amazon-linux-extras install docker
        ;;

        *)  # Amazon Linux 1 uses versions like "2018.3" here, so dont bother enumerating
            yum -y install docker
        ;;
    esac
}

install_jq() {
    case "$distro" in
        alpine)
            apk add --update --no-cache jq
            ;;

        amzn)
            yum -y install jq
            ;;

        *)
            export DEBIAN_FRONTEND=noninteractive
            apt_get_update
            apt-get install -y jq
            ;;
    esac
}

if [ "$(id -u)" != 0 ]; then
    echo "Warning: Docker-in-Earthly needs to be run as root user"
fi

if ! detect_jq; then
    echo "jq is missing. Attempting to install automatically."
    install_jq
fi

if ! detect_dockerd; then
    echo "Docker Engine is missing. Attempting to install automatically."
    install_dockerd
    echo "Docker Engine was missing. It has been installed automatically by Earthly."
    dockerd --version
    echo "For better use of cache, try using the official earthly/dind image for WITH DOCKER."
else
    print_debug "dockerd already installed"
fi

set +u
if [ "$EARTHLY_START_COMPOSE" = "true" ] || [ "$EARTHLY_START_COMPOSE" = "" ]; then
    set -u
    set +e;
    docker_compose="$(detect_docker_compose_cmd)"
    set -e
    if [ -z "$docker_compose" ]; then
        echo "Docker Compose is missing. Attempting to install automatically."
        install_docker_compose

        docker_compose="$(detect_docker_compose_cmd)"
        echo "Docker Compose was missing. It has been installed automatically by Earthly."

        $docker_compose --version
        echo "For better use of cache, try using the official earthly/dind image for WITH DOCKER."
    else
        print_debug "docker-compose already installed"
    fi
else
    print_debug "docker-compose not needed"
fi
