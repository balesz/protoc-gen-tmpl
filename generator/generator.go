package generator

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"text/template"

	"github.com/balesz/protoc-gen-tmpl/data"
	"github.com/balesz/protoc-gen-tmpl/functions"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	rxEnum    = regexp.MustCompile(`.*\{\{.* \.Enum.*\}\}.*`)
	rxFile    = regexp.MustCompile(`.*\{\{.* \.File.*\}\}.*`)
	rxMessage = regexp.MustCompile(`.*\{\{.* \.Message.*\}\}.*`)
	rxService = regexp.MustCompile(`.*\{\{.* \.Service.*\}\}.*`)
)

func New(request *pluginpb.CodeGeneratorRequest) *generator {
	d := data.New(request)
	f := functions.New(d)
	return &generator{data: d, functions: f}
}

type generator struct {
	data      *data.Data
	functions *functions.Functions
}

func (it *generator) Execute() *pluginpb.CodeGeneratorResponse {
	templateFiles := it.data.GetTemplateFiles()
	if len(templateFiles) == 0 {
		return responseError(fmt.Errorf("no template file exists"))
	}

	var output []*pluginpb.CodeGeneratorResponse_File

	if out, err := it.executeStaticTemplates(); err != nil {
		return responseError(err)
	} else if len(out) > 0 {
		output = append(output, out...)
	}

	if out, err := it.executeDynamicTemplates(); err != nil {
		return responseError(err)
	} else if len(out) > 0 {
		output = append(output, out...)
	}

	return &pluginpb.CodeGeneratorResponse{
		File:              output,
		SupportedFeatures: proto.Uint64(uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)),
	}
}

func (it *generator) executeStaticTemplates() ([]*pluginpb.CodeGeneratorResponse_File, error) {
	var templateFiles []string
	for _, it := range it.data.GetTemplateFiles() {
		if !strings.Contains(it, "{{") && !strings.Contains(it, "}}") {
			templateFiles = append(templateFiles, it)
		}
	}
	data := map[string]interface{}{"Files": it.data.Files()}
	return it.generate(templateFiles, data)
}

func (it *generator) executeDynamicTemplates() ([]*pluginpb.CodeGeneratorResponse_File, error) {
	templateFiles := make(map[string][]string)
	for _, it := range it.data.GetTemplateFiles() {
		if rxEnum.MatchString(it) {
			templateFiles["enum"] = append(templateFiles["enum"], it)
		} else if rxMessage.MatchString(it) {
			templateFiles["message"] = append(templateFiles["message"], it)
		} else if rxService.MatchString(it) {
			templateFiles["service"] = append(templateFiles["service"], it)
		} else if rxFile.MatchString(it) {
			templateFiles["file"] = append(templateFiles["file"], it)
		}
	}

	var files []protoreflect.FileDescriptor
	it.data.Registry().RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		files = append(files, fd)
		return true
	})

	data := map[string]interface{}{"Files": it.data.Files()}
	var output []*pluginpb.CodeGeneratorResponse_File
	for _, file := range files {
		data["File"] = file
		if out, err := it.generate(templateFiles["file"], data); err != nil {
			return nil, err
		} else {
			output = append(output, out...)
		}

		for i := 0; i < file.Services().Len(); i++ {
			data["Service"] = file.Services().Get(i)
			if out, err := it.generate(templateFiles["service"], data); err != nil {
				return nil, err
			} else {
				output = append(output, out...)
			}
		}

		for i := 0; i < file.Messages().Len(); i++ {
			data["Message"] = file.Messages().Get(i)
			if out, err := it.generate(templateFiles["message"], data); err != nil {
				return nil, err
			} else {
				output = append(output, out...)
			}
		}

		for i := 0; i < file.Enums().Len(); i++ {
			data["Enum"] = file.Enums().Get(i)
			if out, err := it.generate(templateFiles["enum"], data); err != nil {
				return nil, err
			} else {
				output = append(output, out...)
			}
		}
	}
	return output, nil
}

func (it *generator) generate(files []string, data map[string]interface{}) ([]*pluginpb.CodeGeneratorResponse_File, error) {
	var output []*pluginpb.CodeGeneratorResponse_File
	for _, templateFile := range files {
		templateStr, err := readFile(templateFile)
		if err != nil {
			log.Printf("[%v] %v", templateFile, err)
			continue
		}

		var content string
		buf := new(bytes.Buffer)
		if tmpl, err := template.New(templateFile).Funcs(it.functions.Map()).Parse(templateStr); err != nil {
			return nil, err
		} else if err := tmpl.Execute(buf, data); err != nil {
			if message, ok := it.functions.LookupExit(err); ok {
				log.Printf("[%v] %v", templateFile, message)
				continue
			} else if message, ok := it.functions.LookupFail(err); ok {
				return nil, fmt.Errorf("[%v] %v", templateFile, message)
			}
			return nil, err
		} else if buf.Len() == 0 {
			continue
		} else {
			content = buf.String()
		}

		fileName := strings.TrimSuffix(
			strings.TrimPrefix(templateFile, it.data.TemplateDir()+"/"), ".tmpl")

		buf.Reset()
		if tmpl, err := template.New(fileName).Funcs(it.functions.Map()).Parse(fileName); err != nil {
			return nil, err
		} else if err = tmpl.Execute(buf, data); err != nil {
			return nil, err
		} else {
			fileName = buf.String()
		}

		output = append(output, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(fileName),
			Content: proto.String(content),
		})
	}
	return output, nil
}

func readFile(path string) (string, error) {
	var content string
	if input, err := ioutil.ReadFile(path); err != nil {
		return "", err
	} else if len(input) == 0 {
		return "", errors.New("file is empty")
	} else if content = strings.TrimSpace(string(input)); content == "" {
		return "", errors.New("file is empty")
	}
	return content, nil
}

func responseError(err error) *pluginpb.CodeGeneratorResponse {
	return &pluginpb.CodeGeneratorResponse{Error: proto.String(err.Error())}
}
