#!/usr/bin/env bash

OUT="pb"
CLIENT_PATH="proto"
GO_SRC="${GOPATH}/src"
GOGOPROTO_ROOT="$(GO111MODULE=on go list -f '{{ .Dir }}' -m github.com/gogo/protobuf)"
GOGOPROTO_PATH="${GOGOPROTO_ROOT}:${GOGOPROTO_ROOT}/protobuf"
VALIDATE_ROOT="$(GO111MODULE=on go list -f '{{ .Dir }}' -m github.com/envoyproxy/protoc-gen-validate)"

REPLACEMENT=Mgogoproto/gogo.proto=github.com/gogo/protobuf/gogoproto,\
Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types

files=$(find "${CLIENT_PATH}" -type f -name "*.proto" -print | cut -d '/' -f2-)
for file in ${files}; do
    protoc -I="${CLIENT_PATH}" \
            -I="${GO_SRC}" \
            -I="${GOGOPROTO_PATH}" \
            -I="${VALIDATE_ROOT}" \
            --gogofast_out=${REPLACEMENT}:${OUT} \
            --validate_out="lang=gogo:${OUT}" \
            "${file}"
done
