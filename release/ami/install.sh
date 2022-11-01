#!/bin/bash

set -xe

export DEBIAN_FRONTEND=noninteractive

sudo apt-get update
sudo apt-get -y install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common \
    lsb-release
 
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
 
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://pkg.earthly.dev/earthly.pgp | sudo gpg --dearmor -o /etc/apt/keyrings/earthly-archive-keyring.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/earthly-archive-keyring.gpg] https://pkg.earthly.dev/deb \
  stable main" | sudo tee /etc/apt/sources.list.d/earthly.list > /dev/null

sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io earthly="$EARTHLY_VERSION"

docker -v
earthly -v

sudo usermod -a -G docker ubuntu
