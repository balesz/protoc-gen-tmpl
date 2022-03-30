# protoc-gen-tmpl

Based on the interfaces of the protoreflect package: <https://pkg.go.dev/google.golang.org/protobuf/reflect/protoreflect>

## Options

## Config

The config name should be `protoc-gen-tmpl.yaml`.

### Dart Example

```yaml
exclude:
  - lib/logic/logic_old.dart.tmpl
output:
  - name: model.dart.tmpl
    path: "lib/model/{{.Message.Name|ToSnake}}.dart"
  - name: service.dart.tmpl
    path: "lib/service/{{.Service.Name|ToSnake}}.dart"
types:
  - enum: "{{.Name|ToCamel}}"
  - message: "{{.Name|ToCamel}}"
  - map: "Map<{{TKey}}, {{TValue}}>"
  - repeated: "List<{{T}}>"
  - scalar:
    double: "double"
    float: "double"
    int32: "int"
    int64: "int"
    uint32: "int"
    uint64: "int"
    sint32: "int"
    sint64: "int"
    fixed32: "int"
    fixed64: "int"
    sfixed32: "int"
    sfixed64: "int"
    bool: "bool"
    string: "String"
    bytes: "Uint8List"
```

## Functions

### Nil

Nil function returns nil.

### Exit(message)

Exit function stop the handling of the template file, and print the message to the output.

### Fail(message)

Fail function abort the whole compilation process and print the message to the output.

### ToList(object)

ToList function try to cast the given object to an iterable list.

### Get(key)

Get function retrieves the value for the specified key.

### Set(key, value)

Set function sets the value for the specified key.
