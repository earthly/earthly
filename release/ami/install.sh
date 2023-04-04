#!/bin/bash

set -xe

sudo yum-config-manager --add-repo https://pkg.earthly.dev/earthly.repo
sudo yum update -y
sudo amazon-linux-extras install docker

sudo yum install -y earthly

docker -v
earthly -v

sudo usermod -a -G docker ec2-user
