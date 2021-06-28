#!/bin/bash
ext=$?
path=(go list -f '{{.Target}}')
echo "$ext"
if [[ $ext -ne 0]] then
    exit $ext
fi
export PATH=$PATH:path
go install

