#!/bin/bash
path=(go list -f '{{.Target}}')
export PATH=$PATH:path
go install

