#!/bin/sh
 
set -e
 
echo 'Cleaning up after bootstrapping...'
sudo yum clean all
sudo rm -rf /tmp/*
cat /dev/null > ~/.bash_history