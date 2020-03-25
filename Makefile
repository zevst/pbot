all: install gen_proto

.PHONY: install
install:
	go get -u github.com/gogo/protobuf/protoc-gen-gogofaster

.PHONY: gen_proto
gen_proto:
	rm -rf pb/* && bash genproto.sh