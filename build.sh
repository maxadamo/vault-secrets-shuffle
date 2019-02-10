#!/bin/bash
if ! which upx &>/dev/null; then
    echo "please download upx here https://github.com/upx/upx/releases"
    echo "and store the executable in your PATH"
    exit
fi
export GOPATH=${HOME}/vault-secrets-shuffle
REPO_PATH=${GOPATH}/src/github.com/maxadamo/vault-secrets-shuffle

go get github.com/maxadamo/vault-secrets-shuffle
rm -f github.com/bin/vault-secrets-shuffle
cd ${GOPATH}/src/github.com/maxadamo/vault-secrets-shuffle
go build -ldflags "-s -w"
upx --brute ${REPO_PATH}/vault-secrets-shuffle

echo -e "\nthe binary was compiled and it is avilable as:\n - ${REPO_PATH}/vault-secrets-shuffle"
