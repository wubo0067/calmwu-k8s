.PHONY: all

all : proto

proto:
	for d in srv-upper srv-split sci_auth; do \
		for f in $$d/*.proto; do \
			protoc --proto_path=.:${GOPATH}/src:../../../ --go_out=. --micro_out=. $$f; \
			echo compiled: $$f; \
		done \
	done
