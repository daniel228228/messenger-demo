module_name := messenger.api
major_ver := 1

go_change_major_ver_script := scripts/go_change_major_ver.sh
go_move_compiled_files_script := scripts/go_move_compiled_files.sh

proto_files := $(shell find proto -iname "*.proto")

go_generated := go/api

all: clean gen

clean: clean_go_compiled_proto

gen: go

go:
	-protoc $(proto_files) \
	-I include/googleapis \
	--proto_path=proto \
	--go_out=./go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./go \
	--go-grpc_opt=paths=source_relative \
	--go-grpc-gateway_out=./go \
	--go-grpc-gateway_opt=paths=source_relative \
	--experimental_allow_proto3_optional
	$(MAKE) go_move_compiled_files move_dir="$(go_generated)"
	$(MAKE) go_change_major_ver change_dir="$(go_generated)"

go_move_compiled_files:
	-$(shell $(go_move_compiled_files_script) "$(move_dir)")

go_change_major_ver:
	-$(shell $(go_change_major_ver_script) $(module_name) $(major_ver) "$(change_dir)")

clean_go_compiled_proto:
	rm -rf $(go_generated)/*

.PHONY: go
