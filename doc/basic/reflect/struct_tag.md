# 结构体字段标签（Struct field Tag）

格式：

```
`key1:"value1" key2:"value2"`
```

结构体字段标签由一个或多个键值对组成。键与值使用冒号分隔，值用双引号括起来。键值对之间使用一个空格分隔。

在 JSON 序列化 和 ORM框架中很常见。

结构体字段标签定位跟Java中的注解有点类似，都需要通过反射取值，然后通过标签处理器/注解处理器处理。另外在其他语言貌似也见到过这种写法。

结构体对象：

```go
// A StructField describes a single field in a struct.
type StructField struct {
    // Name is the field name.
    Name string

    // PkgPath is the package path that qualifies a lower case (unexported)
    // field name. It is empty for upper case (exported) field names.
    // See https://golang.org/ref/spec#Uniqueness_of_identifiers
    PkgPath string

    Type      Type      // field type
    Tag       StructTag // field tag string，字段标签字符串
    Offset    uintptr   // offset within struct, in bytes
    Index     []int     // index sequence for Type.FieldByIndex
    Anonymous bool      // is an embedded field
}

//就是个string类型
type StructTag string
```

获取结构体字段标签方法：由上面知道获取filed对象实例就能获取到字段的标签。

```go
t2 := reflect.TypeOf(user)
for i := 0; i < t2.NumField(); i++ {
	field := t2.Field(i)
	tag := field.Tag.Get(tagValidate)
}
```





