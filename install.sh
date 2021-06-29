#!/bin/bash
echo "Cheking requirements..."
go version
if [ $? -eq 0 ]
then
    echo "Go is installed"
else
    echo "ERROR: Go does not seem to be installed."
    echo "Please download Go using your package manager or over https://golang.org/"
    exit 1
fi

path=(go list -f '{{.Target}}')
export PATH=$PATH:path


if [ -x "$(go install)"]
then
    echo "Aura was installed succesfully"
else
    echo "ERROR: it was not possible to install aura :("
    exit 1
fi

