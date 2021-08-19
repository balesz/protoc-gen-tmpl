package main

import (
	"io/ioutil"
	"log"
	"os"

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

	generator := NewGenerator(&req)
	response := generator.Execute()

	if out, err := proto.Marshal(response); err != nil {
		log.Fatal(err)
	} else if _, err := os.Stdout.Write(out); err != nil {
		log.Fatal(err)
	}
}
