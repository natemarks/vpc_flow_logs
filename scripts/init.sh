#!/usr/bin/env bash
set -Eeuoxv pipefail
# shellcheck disable=SC1083
go mod init github.com/natemarks/vpc_flow_logs
git init .
# shellcheck disable=SC1083
git config user.name "Nate Marks"
# shellcheck disable=SC1083
git config user.email "npmarks@gmail.com"
git add -A
git commit -am 'initial'
git branch -m master main
