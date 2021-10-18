#!/bin/sh

set -eu

distro=$(sed -n -e 's/^ID="\?\([^\"]*\)"\?/\1/p' /etc/os-release)

detect_dockerd() {
    set +e
    command -v dockerd
    has_d="$?"
    set -e
    return "$has_d"
}

detect_docker_compose() {
    set +e
    command -v docker-compose
    has_dc="$?"
    set -e
    return "$has_dc"
}

print_debug() {
    set +u
    if [ "$EARTHLY_DEBUG" = "true" ] ; then
        echo "$@"
    fi
    set -u
}

install_docker_compose() {
    case "$distro" in
        alpine)
            apk add --update --no-cache docker-compose
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
            apk add --update --no-cache docker
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

install_dockerd_debian_like() {
    export DEBIAN_FRONTEND=noninteractive
    apt-get remove -y docker docker-engine docker.io containerd runc || true
    apt-get update
    apt-get install -y \
        apt-transport-https \
        ca-certificates \
        curl \
        gnupg-agent \
        software-properties-common
    curl -fsSL "https://download.docker.com/linux/$distro/gpg" | apt-key add -
    add-apt-repository \
        "deb [arch=$(dpkg --print-architecture)] https://download.docker.com/linux/$distro \
        $(lsb_release -cs) \
        stable"
    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io
}

install_dockerd_amazon() {
    version=$(sed -n -e 's/^VERSION="\?\([^\"]*\)"\?/\1/p' /etc/os-release)
    case "$version" in
        2)
            yes | amazon-linux-extras install docker
        ;;

        *)  # Amazon Linux 1 uses versions like "2018.3" here, so dont bother enumerating
            yum -y install docker
        ;;
    esac
}

if [ "$(id -u)" != 0 ]; then
    echo "Warning: Docker-in-Earthly needs to be run as root user"
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
    if ! detect_docker_compose; then
        echo "Docker Compose is missing. Attempting to install automatically."
        install_docker_compose
        echo "Docker Compose was missing. It has been installed automatically by Earthly."
        docker-compose --version
        echo "For better use of cache, try using the official earthly/dind image for WITH DOCKER."
    else
        print_debug "docker-compose already installed"
    fi
else
    print_debug "docker-compose not needed"
fi
