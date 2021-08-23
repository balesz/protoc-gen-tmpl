package functions

import (
	"errors"
	"reflect"
)

func (it *Functions) nilFunc() interface{} {
	return nil
}

func (it *Functions) exitFunc(message string) (string, error) {
	return "", errors.New(message)
}

func (it *Functions) failFunc(message string) (string, error) {
	return "", errors.New(message)
}

func (it *Functions) toListFunc(list interface{}) interface{} {
	var result []interface{}
	listValue := reflect.ValueOf(list)
	if listValue.IsZero() {
		return list
	}
	lenFunc := reflect.ValueOf(list).MethodByName("Len")
	getFunc := reflect.ValueOf(list).MethodByName("Get")
	if lenFunc.IsZero() || getFunc.IsZero() {
		return list
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
	return result
}
