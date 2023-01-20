PROTOC=protoc

PROTO_DIR = ./pb

PROTO_GO_FILES:
	$(PROTOC) -I=$(PROTO_DIR) --go-grpc_out=. $(PROTO_DIR)/*.proto
	$(PROTOC) -I=$(PROTO_DIR) --go_out=./ $(PROTO_DIR)/*.proto

.PHONY: pb
pb: PROTO_GO_FILES

.PHONY: clean
clean:
	rm pb/*.pb.go