package interfaces

import (
	"encoding/json"
	"log"
	"testing"
)

type Person struct {
	Name string
	Age  uint8
}

func (p *Person) Sing() {
	log.Println("I can sing")
}

func (p *Person) String() (string, error) {
	bs, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

type Action interface {
	Sing()
}

type EmptyInterface interface{}

// TestTypeConvertToInterface empty interface 可以接受任何类型的参数，但是不能把某类型参数强转成 empty interface
func TestTypeConvertToInterface(t *testing.T) {
	arvin := &Person{
		Name: "Arvin",
		Age:  18,
	}

	var obj interface{} = arvin
	log.Println("obj: ", obj)

	//1 empty interface 可以接受任何类型的参数，但不能把某类型参数强转成 empty interface
	//语法检查时就会报：Invalid type assertion: arvin.(interface{}) (non-interface type *Person on the left)
	//obj, ok := arvin.(interface{})

	persons := make([]*Person, 2)
	persons[0] = arvin
	persons[1] = arvin
	//2 empty interface 可以接受任何类型的参数，但是 `interface{}` 类型的 slice 不可以接受任何类型的 slice
	//var objs []interface{}
	//objs = persons
}
