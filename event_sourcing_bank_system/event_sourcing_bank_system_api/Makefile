PROTO_DIR := proto

server:
	go run cmd/main.go
	
protos:
	@go install golang.org/x/tools/cmd/goimports@v0.1.12
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	@protoc \
	  --experimental_allow_proto3_optional \
	  --proto_path=$(PROTO_DIR) \
	  --go_out=paths=source_relative:$(PROTO_DIR) \
	  --go-grpc_out=require_unimplemented_servers=false,paths=source_relative:$(PROTO_DIR) \
	  $(PROTO_DIR)/payment/*.proto
	@$(GOIMPORTS) -w $(PROTO_DIR)
.PHONY: protos

.PHONY: server protos