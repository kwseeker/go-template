package reflect

import (
	"fmt"
	"reflect"
	"testing"
)

// 先通过反射获取字段对象（StructField），然后获取 Tag字段(string类型)，get方法从string Tag中查找key对应的值
func TestGetTag(t *testing.T) {
	user := User{
		Id:    1,
		Name:  "John Doe",
		Email: "john@example",
	}

	// TypeOf returns the reflection Type that represents the dynamic type of variable.
	// If variable is a nil interface value, TypeOf returns nil.
	t2 := reflect.TypeOf(user)

	//Get the type and kind of our user variable
	fmt.Println("Type: ", t2.Name())
	fmt.Println("Kind: ", t2.Kind())

	for i := 0; i < t2.NumField(); i++ {
		field := t2.Field(i)
		//获取字段的标签
		tag := field.Tag.Get(tagValidate)
		fmt.Printf("%d. %v(%v), tag:'%v'\n", i+1, field.Name, field.Type.Name(), tag)
	}
}
