#!/bin/sh
 
set -e
 
echo 'Cleaning up after bootstrapping...'
sudo apt-get -y autoremove
sudo apt-get -y clean
sudo rm -rf /tmp/*
cat /dev/null > ~/.bash_history
history -c
exit