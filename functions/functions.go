package functions

import (
	"errors"
	"reflect"
	"regexp"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/balesz/protoc-gen-tmpl/data"
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func New(data *data.Data) *Functions {
	return &Functions{data: data, store: make(map[string]interface{})}
}

type Functions struct {
	data    *data.Data
	funcMap template.FuncMap
	store   map[string]interface{}
}

func (it *Functions) Map() template.FuncMap {
	if len(it.funcMap) != 0 {
		return it.funcMap
	}
	it.funcMap = template.FuncMap{
		"Nil":    it.nilFunc,
		"Exit":   it.exitFunc,
		"Fail":   it.failFunc,
		"ToList": it.toListFunc,

		"Set": it.setFunc,
		"Get": it.getFunc,

		"ToCamel":              strcase.ToCamel,
		"ToDelimited":          strcase.ToDelimited,
		"ToKebab":              strcase.ToKebab,
		"ToLowerCamel":         strcase.ToLowerCamel,
		"ToScreamingDelimited": strcase.ToScreamingDelimited,
		"ToScreamingKebab":     strcase.ToScreamingKebab,
		"ToScreamingSnake":     strcase.ToScreamingSnake,
		"ToSnake":              strcase.ToSnake,
		"ToSnakeWithIgnore":    strcase.ToSnakeWithIgnore,

		"FindDescriptorByName": it.findDescriptorByNameFunc,
		"FindEnumByName":       it.findEnumByNameFunc,
		"FindMessageByName":    it.findMessageByNameFunc,
		"FindFileByPath":       it.findFileByPathFunc,
		"FindServiceByName":    it.findServiceByNameFunc,
	}
	for k, v := range sprig.TxtFuncMap() {
		it.funcMap[k] = v
	}
	return it.funcMap
}

func (it *Functions) LookupExit(err error) (string, bool) {
	rxExit := regexp.MustCompile(".* error calling (Exit: .*)")
	if rxExit.MatchString(err.Error()) {
		matches := rxExit.FindStringSubmatch(err.Error())
		return matches[1], true
	}
	return "", false
}

func (it *Functions) LookupFail(err error) (string, bool) {
	rxFail := regexp.MustCompile(".* error calling (Fail: .*)")
	if rxFail.MatchString(err.Error()) {
		matches := rxFail.FindStringSubmatch(err.Error())
		return matches[1], true
	}
	return "", false
}

func (it *Functions) ResetStore() {
	it.store = make(map[string]interface{})
}

func (it *Functions) nilFunc() interface{} {
	return nil
}

func (it *Functions) exitFunc(message string) (string, error) {
	return "", errors.New(message)
}

func (it *Functions) failFunc(message string) (string, error) {
	return "", errors.New(message)
}

func (it *Functions) toListFunc(list interface{}) []interface{} {
	var result []interface{}
	lenFunc := reflect.ValueOf(list).MethodByName("Len")
	getFunc := reflect.ValueOf(list).MethodByName("Get")
	if lenFunc.IsZero() || getFunc.IsZero() {
		return result
	}
	length := lenFunc.Call([]reflect.Value{})[0]
	for i := 0; i < int(length.Int()); i++ {
		value := getFunc.Call([]reflect.Value{reflect.ValueOf(i)})[0]
		result = append(result, value.Interface())
	}
	return result
}

func (it *Functions) setFunc(key string, val interface{}) interface{} {
	it.store[key] = val
	return ""
}

func (it *Functions) getFunc(key string) interface{} {
	result := it.store[key]
	delete(it.store, key)
	return result
}

func (it *Functions) findDescriptorByNameFunc(name string) protoreflect.Descriptor {
	return it.data.FindDescriptorByName(name)
}

func (it *Functions) findEnumByNameFunc(name string) protoreflect.EnumDescriptor {
	if result, ok := it.data.FindDescriptorByName(name).(protoreflect.EnumDescriptor); ok {
		return result
	} else {
		return nil
	}
}

func (it *Functions) findFileByPathFunc(path string) protoreflect.FileDescriptor {
	return it.data.FindFileByPath(path)
}

func (it *Functions) findMessageByNameFunc(name string) protoreflect.MessageDescriptor {
	if result, ok := it.data.FindDescriptorByName(name).(protoreflect.MessageDescriptor); ok {
		return result
	} else {
		return nil
	}
}

func (it *Functions) findServiceByNameFunc(name string) protoreflect.ServiceDescriptor {
	if result, ok := it.data.FindDescriptorByName(name).(protoreflect.ServiceDescriptor); ok {
		return result
	} else {
		return nil
	}
}
