# protoc-gen-tmpl

Based on the interfaces of the protoreflect package: <https://pkg.go.dev/google.golang.org/protobuf/reflect/protoreflect>

## Config

### Dart Example

```yaml
exclude:
  - lib/logic/logic_old.dart.tmpl
output:
  - name: lib/model/model.dart.tmpl
    path: "lib/model/{{.Message.Name|ToSnake}}.dart"
  - name: lib/logic/logic.dart.tmpl
    path: "lib/logic/{{.Service.Name|ToSnake}}.dart"
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
