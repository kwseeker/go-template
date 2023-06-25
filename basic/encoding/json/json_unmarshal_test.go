package json

import (
	"encoding/json"
	"fmt"
	"testing"
)

/*
有同事初学golang，问一个请求Handler适配多个来源，请求的json结构大致相同但是有小的差别，应该怎么定义结构体并处理
有两种方式：
1 在请求的时候加上来源信息，然后将动态的这个字段定义成泛型，请求进来的时候根据来源信息选择泛型实际类型，进行转换；
2 将这个动态字段定义成通用类型，请求进来的时候先不转这个字段，到实际需要读取这个字段的内容时再转成对应的这个类型。
*/

type Person struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}

// 方式2
//type Father struct {
//	Person
//	//假设Child是动态类型，可能是Person或其他类型
//	//Child  string `json:"child"`
//	Child interface{} `json:"child"`
//	//Child map[string]interface{} `json:"child"` // encoding/json 内部使用 map[string]interface{} 作为中间数据类型
//}

type Animal struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  string `json:"sex"`
}

type Children interface {
	Person | Animal
}

// 方式1：泛型
type Father[T Children] struct {
	Person
	//假设Child是动态类型，可能是Person或Animal
	Child T `json:"child"`
}

func TestJsonUnmarshal(t *testing.T) {
	jsonStr := `{"name": "Tom", "age": 20, "gender": "male", "child": {"name": "Tom", "age": 20, "gender": "male"}}`
	jsonAnimalStr := `{"name": "Tom", "age": 20, "gender": "male", "child": {"name": "Tom", "age": 20, "sex": "male"}}`
	var p Father[Person]
	var p2 Father[Animal]
	err := json.Unmarshal([]byte(jsonStr), &p)
	if err != nil {
		fmt.Println("解析 JSON 数据失败：", err)
	}
	err = json.Unmarshal([]byte(jsonAnimalStr), &p2)
	if err != nil {
		fmt.Println("解析 JSON 数据失败：", err)
	}
	fmt.Printf("反序列化后的结构体：%v, %v\n", p, p2)

	//实际需要访问Child数据时，再转成对应类型比如Person
	//childStr, _ := json.Marshal(p.Child)
	//var child Person
	//err = json.Unmarshal(childStr, &child)
	//fmt.Printf("反序列化后的结构体：%v\n", child)
}
