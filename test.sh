#!/bin/sh

strs="123,456,789"

for i in $(echo "$strs" | sed "s/,/ /g")
do
    # call your procedure/other scripts here below
    echo "$i"
done
