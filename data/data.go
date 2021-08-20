package data

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func New(request *pluginpb.CodeGeneratorRequest) *Data {
	registry, err := protodesc.NewFiles(&descriptorpb.FileDescriptorSet{File: request.ProtoFile})
	if err != nil {
		panic(err)
	}
	return &Data{request, registry}
}

type Data struct {
	request  *pluginpb.CodeGeneratorRequest
	registry *protoregistry.Files
}

func (it *Data) Request() *pluginpb.CodeGeneratorRequest {
	return it.request
}

func (it *Data) Registry() *protoregistry.Files {
	return it.registry
}

func (it *Data) Files() []protoreflect.FileDescriptor {
	var files []protoreflect.FileDescriptor
	it.registry.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		files = append(files, fd)
		return true
	})
	return files
}

func (it *Data) FindDescriptorByName(name string) protoreflect.Descriptor {
	result, _ := it.registry.FindDescriptorByName(protoreflect.FullName(name))
	return result
}

func (it *Data) FindFileByPath(path string) protoreflect.FileDescriptor {
	result, _ := it.registry.FindFileByPath(path)
	return result
}

func (it *Data) TemplateDir() string {
	if templateDir := it.request.GetParameter(); templateDir == "" {
		return "template"
	} else {
		return templateDir
	}
}

func (it *Data) GetTemplateFiles() []string {
	var result []string
	filepath.Walk(it.TemplateDir(), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("%v", err)
		} else if info.IsDir() {
			return nil
		} else if !strings.HasSuffix(info.Name(), ".tmpl") {
			log.Printf("invalid template file: %v", path)
		} else {
			result = append(result, path)
		}
		return nil
	})
	return result
}
