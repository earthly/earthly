#!/bin/sh

set -eu

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

for_alpine() {
    if ! detect_dockerd; then
        echo "Docker Engine is missing. Attempting to install automatically."
        apk add --update --no-cache docker
        echo "Docker Engine was missing. It has been installed automatically by Earthly."
        dockerd --version
        echo "For better use of cache, try using the official earthly/dind image for WITH DOCKER."
    else
        print_debug "dockerd already installed"
    fi
    if [ "$EARTHLY_START_COMPOSE" = "true" ]; then
        if ! detect_docker_compose; then
            echo "Docker Compose is missing. Attempting to install automatically."
            apk add --update --no-cache docker-compose
            echo "Docker Compose was missing. It has been installed automatically by Earthly."
            docker-compose --version
            echo "For better use of cache, try using the official earthly/dind image for WITH DOCKER."
        else
            print_debug "docker-compose already installed"
        fi
    else
        print_debug "docker-compose not needed"
    fi
}

for_debian() {
    if ! detect_dockerd; then
        echo "Docker Engine is missing. Attempting to install automatically."
        apt-get update
        apt-get install -y \
            apt-transport-https \
            ca-certificates \
            curl \
            gnupg-agent \
            software-properties-common
        curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
        add-apt-repository \
            "deb [arch=amd64] https://download.docker.com/linux/debian \
            $(lsb_release -cs) \
            stable"
        apt-get update
        apt-get install -y docker-ce docker-ce-cli containerd.io
        echo "Docker Engine was missing. It has been installed automatically by Earthly."
        dockerd --version
        echo "For better use of cache, try using the official earthly/dind image for WITH DOCKER."
    else
        print_debug "dockerd already installed"
    fi

    if [ "$EARTHLY_START_COMPOSE" = "true" ]; then
        if ! detect_docker_compose; then
            echo "Docker Compose is missing. Attempting to install automatically."
            curl -L "https://github.com/docker/compose/releases/download/1.27.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
            chmod +x /usr/local/bin/docker-compose
            echo "Docker Compose was missing. It has been installed automatically by Earthly."
            docker-compose --version
            echo "For better use of cache, try using the official earthly/dind image for WITH DOCKER."
        else
            print_debug "docker-compose already installed"
        fi
    else
        print_debug "docker-compose not needed"
    fi
}

if [ "$(id -u)" != 0 ]; then
    echo "Warning: Docker-in-Earthly needs to be run as root user"
fi

distro=$(sed -n -e 's/^ID=\(.*\)/\1/p' /etc/os-release)
case "$distro" in
    alpine)
        for_alpine
        ;;

    ubuntu)
        for_debian
        ;;

    debian)
        for_debian
        ;;

    *)
        echo "Warning: Distribution $distro not yet supported for Docker-in-Earthly."
        echo "Will attempt to treat like Debian."
        for_debian
        ;;
esac
