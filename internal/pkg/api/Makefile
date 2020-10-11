PROTO_DIR = .
GO_DIR = .

PROTOC = protoc
PROTOTOOL = prototool

PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)
GO_FILES = $(PROTO_FILES:$(PROTO_DIR)/%.proto=$(GO_DIR)/%.pb.go)

$(GO_DIR)/%.pb.go: $(PROTO_DIR)/%.proto FORCE
	$(PROTOC) -I=$(PROTO_DIR) --go_out=plugins=grpc,paths=source_relative:$(GO_DIR) $<

default: build

build: $(GO_FILES)

lint:
	$(PROTOTOOL) lint $(PROTO_DIR)

clean:
	rm $(GO_DIR)/*.pb.go

FORCE:

.PHONY: default build lint clean
