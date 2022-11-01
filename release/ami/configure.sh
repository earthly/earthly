#!/bin/bash

set -xe

sudo systemctl start docker.service
sudo systemctl start containerd.service

sudo systemctl enable docker.service
sudo systemctl enable containerd.service

docker run --rm --privileged aptman/qus -s -- --r -p
earthly bootstrap
