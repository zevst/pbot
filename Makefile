all: install gen_proto

.PHONY: install
install:
	go get -u \
		github.com/gogo/protobuf/protoc-gen-gogofaster \
		github.com/envoyproxy/protoc-gen-validate

.PHONY: gen_proto
gen_proto:
	rm -rf pb/* && bash genproto.sh