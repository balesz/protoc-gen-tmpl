BIN_NAME := protoc-gen-tmpl
INPUT_DIR := ../protos_gamerpro/remux/app
OUTPUT_DIR := test/generated
DESCRIPTOR_OUT := test/descriptor.out

GOOGLE_PROTO_FILES := google/protobuf/any.proto google/protobuf/empty.proto google/protobuf/timestamp.proto
PROTO_FILES := $(notdir $(wildcard $(INPUT_DIR)/*.proto))

clean:
	rm -rf bin $(DESCRIPTOR_OUT) $(OUTPUT_DIR)

descriptor.out: $(INPUT_DIR)/remux.proto
	rm $(DESCRIPTOR_OUT)
	protoc -I=$(INPUT_DIR) -o$(DESCRIPTOR_OUT) \
	$(GOOGLE_PROTO_FILES) \
	$(PROTO_FILES)

protoc-gen-tmpl: main.go
	rm bin/$(BIN_NAME) || true
	go build -o bin/$(BIN_NAME) .

remux.dart: $(BIN_NAME) $(INPUT_DIR)/remux.proto 
	rm -rf $(OUTPUT_DIR); mkdir $(OUTPUT_DIR)
	protoc -I=$(INPUT_DIR) \
	--plugin=$(BIN_NAME)=bin/$(BIN_NAME) \
	--tmpl_out=$(OUTPUT_DIR) --tmpl_opt=test \
	$(GOOGLE_PROTO_FILES) \
	$(PROTO_FILES)


.PHONY: clean descriptor.out protoc-gen-tmpl remux.dart $(BIN_NAME)
