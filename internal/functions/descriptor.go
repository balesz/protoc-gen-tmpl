package functions

import (
	"bytes"
	"strings"
	"text/template"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func (it *Functions) mapType(desc protoreflect.FieldDescriptor) string {
	data := struct {
		protoreflect.FieldDescriptor
		Type, TypeKey, TypeValue string
	}{desc, "unknown", "unknown", "unknown"}

	switch {
	case desc.IsMap():
		data.TypeKey = it.mapType(desc.MapKey())
		data.TypeValue = it.mapType(desc.MapValue())
		buf := new(bytes.Buffer)
		format := it.data.Config().TypeFormat("map")
		if temp, err := template.New("").Parse(format); err != nil {
			return "unknown"
		} else if err := temp.Execute(buf, data); err != nil {
			return "unknown"
		}
		data.Type = buf.String()
	case desc.IsList():
		buf := new(bytes.Buffer)
		format := it.data.Config().TypeFormat(desc.Kind().String())
		if temp, err := template.New("").Parse(format); err != nil {
			return "unknown"
		} else if err := temp.Execute(buf, data); err != nil {
			return "unknown"
		}
		data.Type = buf.String()
		buf = new(bytes.Buffer)
		format = it.data.Config().TypeFormat("list")
		if temp, err := template.New("").Parse(format); err != nil {
			return "unknown"
		} else if err := temp.Execute(buf, data); err != nil {
			return "unknown"
		}
		data.Type = buf.String()
	default:
		buf := new(bytes.Buffer)
		format := it.data.Config().TypeFormat(desc.Kind().String())
		if temp, err := template.New("").Parse(format); err != nil {
			return "unknown"
		} else if err := temp.Execute(buf, data); err != nil {
			return "unknown"
		}
		data.Type = buf.String()
	}

	if desc.HasOptionalKeyword() {
		buf := new(bytes.Buffer)
		format := it.data.Config().TypeFormat("optional")
		if temp, err := template.New("").Parse(format); err != nil {
			return "unknown"
		} else if err := temp.Execute(buf, data); err != nil {
			return "unknown"
		}
		data.Type = buf.String()
	}

	return data.Type
}

func (it *Functions) findFileByPathFunc(path string) protoreflect.FileDescriptor {
	result, _ := it.data.Registry().FindFileByPath(path)
	return result
}

func (it *Functions) findDescriptorByNameFunc(name string) protoreflect.Descriptor {
	result, _ := it.data.Registry().FindDescriptorByName(protoreflect.FullName(name))
	return result
}

func (it *Functions) findServiceByNameFunc(name string) protoreflect.ServiceDescriptor {
	if result, ok := it.findDescriptorByNameFunc(name).(protoreflect.ServiceDescriptor); ok {
		return result
	} else {
		return nil
	}
}

func (it *Functions) findMessageByNameFunc(name string) protoreflect.MessageDescriptor {
	if result, ok := it.findDescriptorByNameFunc(name).(protoreflect.MessageDescriptor); ok {
		return result
	} else {
		return nil
	}
}

func (it *Functions) findEnumByNameFunc(name string) protoreflect.EnumDescriptor {
	if result, ok := it.findDescriptorByNameFunc(name).(protoreflect.EnumDescriptor); ok {
		return result
	} else {
		return nil
	}
}

func (it *Functions) leadingCommentsFunc(desc protoreflect.Descriptor) string {
	srcLoc := desc.ParentFile().SourceLocations().ByDescriptor(desc)
	return srcLoc.LeadingComments
}

func (it *Functions) leadingDetachedCommentsFunc(desc protoreflect.Descriptor) []string {
	srcLoc := desc.ParentFile().SourceLocations().ByDescriptor(desc)
	return srcLoc.LeadingDetachedComments
}

func (it *Functions) trailingCommentsFunc(desc protoreflect.Descriptor) string {
	srcLoc := desc.ParentFile().SourceLocations().ByDescriptor(desc)
	return strings.TrimSpace(srcLoc.TrailingComments)
}

func (it *Functions) optionsFunc(desc protoreflect.Descriptor) interface{} {
	if _, ok := desc.(protoreflect.FileDescriptor); ok {
		return desc.Options().(*descriptorpb.FileOptions)
	} else if _, ok := desc.(protoreflect.MessageDescriptor); ok {
		return desc.Options().(*descriptorpb.MessageOptions)
	} else if _, ok := desc.(protoreflect.FieldDescriptor); ok {
		return desc.Options().(*descriptorpb.FieldOptions)
	} else if _, ok := desc.(protoreflect.OneofDescriptor); ok {
		return desc.Options().(*descriptorpb.OneofOptions)
	} else if _, ok := desc.(protoreflect.EnumDescriptor); ok {
		return desc.Options().(*descriptorpb.EnumOptions)
	} else if _, ok := desc.(protoreflect.EnumValueDescriptor); ok {
		return desc.Options().(*descriptorpb.EnumValueOptions)
	} else if _, ok := desc.(protoreflect.ServiceDescriptor); ok {
		return desc.Options().(*descriptorpb.ServiceOptions)
	} else if _, ok := desc.(protoreflect.MethodDescriptor); ok {
		return desc.Options().(*descriptorpb.MethodOptions)
	}
	return nil
}
