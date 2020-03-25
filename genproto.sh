#!/usr/bin/env bash

OUT="pb"
CLIENT_PATH="proto"
GO_SRC="${GOPATH}/src"
GOGOPROTO_ROOT="$(GO111MODULE=on go list -f '{{ .Dir }}' -m github.com/gogo/protobuf)"
GOGOPROTO_PATH="${GOGOPROTO_ROOT}:${GOGOPROTO_ROOT}/protobuf"

REPLACEMENT=Mgogoproto/gogo.proto=github.com/gogo/protobuf/gogoproto

files=$(find "${CLIENT_PATH}" -type f -name "*.proto" -print | cut -d '/' -f2-)
for file in ${files}; do
    protoc -I="${CLIENT_PATH}" \
            -I="${GO_SRC}" \
            -I="${GOGOPROTO_PATH}" \
            --gogofast_out=${REPLACEMENT}:${OUT} \
            "${file}"
done
