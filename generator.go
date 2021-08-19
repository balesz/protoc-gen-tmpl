package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	rxExit = regexp.MustCompile(".* error calling (exit: .*)")
	rxFail = regexp.MustCompile(".* error calling (fail: .*)")
)

func NewGenerator(request *pluginpb.CodeGeneratorRequest) *generator {
	templateDir := request.GetParameter()
	if templateDir == "" {
		templateDir = "template"
	}

	var files []protoreflect.FileDescriptor
	dset, err := protodesc.NewFiles(&descriptorpb.FileDescriptorSet{File: request.ProtoFile})
	if err != nil {
		panic(err)
	}
	dset.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		files = append(files, fd)
		return true
	})

	return &generator{templateDir, files}
}

type generator struct {
	templateDir string
	files       []protoreflect.FileDescriptor
}

func (it *generator) Execute() *pluginpb.CodeGeneratorResponse {
	templateFiles := getTemplateFiles(it.templateDir)
	if len(templateFiles) == 0 {
		return responseError(fmt.Errorf("no template file exists"))
	}

	var output []*pluginpb.CodeGeneratorResponse_File
	for _, templateFile := range templateFiles {
		var templateStr string
		if input, err := ioutil.ReadFile(templateFile); err != nil {
			return responseError(err)
		} else if len(input) != 0 {
			templateStr = string(input)
		} else if templateStr == "" {
			continue
		}

		buf := new(bytes.Buffer)
		data := map[string]interface{}{"Files": it.files}
		if tmpl, err := template.New(templateFile).Funcs(FunctionMap).Parse(templateStr); err != nil {
			return responseError(fmt.Errorf("template parse error: %v", err))
		} else if err := tmpl.Execute(buf, data); err != nil {
			if rxExit.MatchString(err.Error()) {
				matches := rxExit.FindStringSubmatch(err.Error())
				log.Printf("[%v] %v", templateFile, matches[1])
				continue
			} else if rxFail.MatchString(err.Error()) {
				matches := rxFail.FindStringSubmatch(err.Error())
				return responseError(fmt.Errorf("[%v] %v", templateFile, matches[1]))
			}
			return responseError(err)
		} else if buf.Len() == 0 {
			continue
		}

		fileName := strings.TrimSuffix(
			strings.TrimPrefix(templateFile, it.templateDir+"/"), ".tmpl")

		output = append(output, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(fileName),
			Content: proto.String(buf.String()),
		})
	}

	return &pluginpb.CodeGeneratorResponse{
		File:              output,
		SupportedFeatures: proto.Uint64(uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)),
	}
}

func getTemplateFiles(root string) []string {
	var result []string
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
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

func responseError(err error) *pluginpb.CodeGeneratorResponse {
	return &pluginpb.CodeGeneratorResponse{Error: proto.String(err.Error())}
}
