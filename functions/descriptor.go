package functions

import "google.golang.org/protobuf/reflect/protoreflect"

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
	return srcLoc.TrailingComments
}
