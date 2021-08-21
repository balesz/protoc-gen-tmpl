package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/balesz/protoc-gen-tmpl/generator"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	var req pluginpb.CodeGeneratorRequest
	if in, err := ioutil.ReadAll(os.Stdin); err != nil {
		log.Fatal(err)
	} else if err := proto.Unmarshal(in, &req); err != nil {
		log.Fatal(err)
	}

	var response *pluginpb.CodeGeneratorResponse
	if files, err := generator.New(&req).Execute(); err != nil {
		response = &pluginpb.CodeGeneratorResponse{
			Error: proto.String(err.Error()),
		}
	} else {
		response = &pluginpb.CodeGeneratorResponse{
			File:              files,
			SupportedFeatures: proto.Uint64(uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)),
		}
	}

	if out, err := proto.Marshal(response); err != nil {
		log.Fatal(err)
	} else if _, err := os.Stdout.Write(out); err != nil {
		log.Fatal(err)
	}
}
