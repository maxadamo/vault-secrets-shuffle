#!/bin/bash
#if ! which upx &>/dev/null; then
#    echo "please download upx here https://github.com/upx/upx/releases"
#    echo "and store the executable within your \$PATH"
#    exit
#fi
BIN_NAME=vault-secrets-shuffle
PATH=$PATH:$(go env GOPATH)/bin
GOPATH=$(go env GOPATH)
export BIN_NAME PATH GOPATH

LATEST_TAG=$(git describe --tags $(git rev-list --tags --max-count=1))
PROG_VERSION=$(echo $LATEST_TAG | sed -e 's/^v//')
BUILDTIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

rm -rf ${GOPATH}/src/github.com/maxadamo/${BIN_NAME}
go get -ldflags "-s -w -X main.appVersion=${PROG_VERSION} -X main.buildTime=${BUILDTIME}" github.com/maxadamo/${BIN_NAME}
# upx --brute ${GOPATH}/bin/${BIN_NAME}

if [ $? -gt 0 ]; then
  echo -e "\nthere was an error while compiling the code\n"
  exit
fi

echo -e "\nthe binary was compiled and it is avilable as:\n - ${GOPATH}/bin/${BIN_NAME}\n"
