package data

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
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
	cfg, _ := InitConfig(request.GetParameter())
	return &Data{request: request, registry: registry, config: cfg}
}

type Data struct {
	request       *pluginpb.CodeGeneratorRequest
	registry      *protoregistry.Files
	config        *Config
	protoFiles    []protoreflect.FileDescriptor
	templateFiles map[string]string
}

func (it *Data) Request() *pluginpb.CodeGeneratorRequest {
	return it.request
}

func (it *Data) Registry() *protoregistry.Files {
	return it.registry
}

func (it *Data) Files() []protoreflect.FileDescriptor {
	if len(it.protoFiles) != 0 {
		return it.protoFiles
	}
	it.registry.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		it.protoFiles = append(it.protoFiles, fd)
		return true
	})
	sort.SliceStable(it.protoFiles, func(i, j int) bool {
		name1 := fmt.Sprintf("%v", it.protoFiles[i].Path())
		name2 := fmt.Sprintf("%v", it.protoFiles[j].Path())
		ord := []string{name1, name2}
		sort.Strings(ord)
		return ord[0] == name1
	})
	return it.protoFiles
}

func (it *Data) TemplateDir() string {
	if templateDir := it.request.GetParameter(); templateDir == "" {
		return "template"
	} else {
		return templateDir
	}
}

func (it *Data) TemplateFiles() map[string]string {
	if it.templateFiles != nil {
		return it.templateFiles
	}
	it.templateFiles = make(map[string]string)
	filepath.Walk(it.TemplateDir(), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("%v", err)
		} else if info.IsDir() {
			return nil
		} else if !strings.HasSuffix(info.Name(), ".tmpl") {
			return nil
		} else if it.config.IsExcluded(path) {
			log.Printf("template file excluded: %v", path)
			return nil
		}
		if out := it.config.OutputByName(path); out == "" {
			it.templateFiles[path] = filepath.Join(strings.Split(path, string(filepath.Separator))[1:]...)
		} else {
			it.templateFiles[path] = out
		}
		it.templateFiles[path] = strings.TrimSuffix(it.templateFiles[path], ".tmpl")
		return nil
	})
	return it.templateFiles
}
