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
	rxEnum    = regexp.MustCompile(`.*\{\{.*\.Enum.*\}\}.*`)
	rxFile    = regexp.MustCompile(`.*\{\{.*\.File.*\}\}.*`)
	rxMessage = regexp.MustCompile(`.*\{\{.*\.Message.*\}\}.*`)
	rxService = regexp.MustCompile(`.*\{\{.*\.Service.*\}\}.*`)
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

func (gen *generator) executeStaticTemplates() ([]*pluginpb.CodeGeneratorResponse_File, error) {
	templateFiles := make(map[string]string)
	for tmpl, out := range gen.data.GetTemplateFiles() {
		if !strings.Contains(out, "{{") && !strings.Contains(out, "}}") {
			templateFiles[tmpl] = out
		}
	}
	data := map[string]interface{}{"Files": gen.data.Files()}
	return gen.generate(templateFiles, data)
}

func (gen *generator) executeDynamicTemplates() ([]*pluginpb.CodeGeneratorResponse_File, error) {
	templateFiles := make(map[string]map[string]string)
	for tmpl, out := range gen.data.GetTemplateFiles() {
		if rxEnum.MatchString(out) {
			if templateFiles["enum"] == nil {
				templateFiles["enum"] = make(map[string]string)
			}
			templateFiles["enum"][tmpl] = out
		} else if rxMessage.MatchString(out) {
			if templateFiles["message"] == nil {
				templateFiles["message"] = make(map[string]string)
			}
			templateFiles["message"][tmpl] = out
		} else if rxService.MatchString(out) {
			if templateFiles["service"] == nil {
				templateFiles["service"] = make(map[string]string)
			}
			templateFiles["service"][tmpl] = out
		} else if rxFile.MatchString(out) {
			if templateFiles["file"] == nil {
				templateFiles["file"] = make(map[string]string)
			}
			templateFiles["file"][tmpl] = out
		}
	}

	var files []protoreflect.FileDescriptor
	gen.data.Registry().RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		files = append(files, fd)
		return true
	})

	data := map[string]interface{}{"Files": gen.data.Files()}
	var output []*pluginpb.CodeGeneratorResponse_File
	for _, file := range files {
		data["File"] = file
		if out, err := gen.generate(templateFiles["file"], data); err != nil {
			return nil, err
		} else {
			output = append(output, out...)
		}

		for i := 0; i < file.Services().Len(); i++ {
			data["Service"] = file.Services().Get(i)
			if out, err := gen.generate(templateFiles["service"], data); err != nil {
				return nil, err
			} else {
				output = append(output, out...)
			}
		}

		for i := 0; i < file.Messages().Len(); i++ {
			data["Message"] = file.Messages().Get(i)
			if out, err := gen.generate(templateFiles["message"], data); err != nil {
				return nil, err
			} else {
				output = append(output, out...)
			}
		}

		for i := 0; i < file.Enums().Len(); i++ {
			data["Enum"] = file.Enums().Get(i)
			if out, err := gen.generate(templateFiles["enum"], data); err != nil {
				return nil, err
			} else {
				output = append(output, out...)
			}
		}
	}
	return output, nil
}

func (gen *generator) generate(files map[string]string, data map[string]interface{}) ([]*pluginpb.CodeGeneratorResponse_File, error) {
	var output []*pluginpb.CodeGeneratorResponse_File
	for templateFile, outputFile := range files {
		templateStr, err := readFile(templateFile)
		if err != nil {
			log.Printf("[%v] %v", templateFile, err)
			continue
		}

		gen.functions.ResetStore()

		var content string
		buf := new(bytes.Buffer)
		if tmpl, err := template.New(templateFile).Funcs(gen.functions.Map()).Parse(templateStr); err != nil {
			return nil, err
		} else if err := tmpl.Execute(buf, data); err != nil {
			if message, ok := gen.functions.LookupExit(err); ok {
				log.Printf("[%v] %v", templateFile, message)
				continue
			} else if message, ok := gen.functions.LookupFail(err); ok {
				return nil, fmt.Errorf("[%v] %v", templateFile, message)
			}
			return nil, err
		} else if buf.Len() == 0 {
			continue
		} else {
			content = buf.String()
		}

		buf.Reset()
		if tmpl, err := template.New(outputFile).Funcs(gen.functions.Map()).Parse(outputFile); err != nil {
			return nil, err
		} else if err = tmpl.Execute(buf, data); err != nil {
			return nil, err
		} else {
			outputFile = buf.String()
		}

		output = append(output, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(outputFile),
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
