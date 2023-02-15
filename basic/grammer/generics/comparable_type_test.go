package generics

import (
	"fmt"
	"testing"
)

//两个变量如果要用“==”来比较，需要满足如下条件：
//	1) 必须是同类型。如果是两个接口则其中一个接口必须定义了另一个接口的全部方法；如果是结构体则必须是同一个命名结构体的两个实例或者是两个相同（包括字段顺序相同）的匿名结构体的实例。
//	2) 不能是func, map, slice
//	3) 如果是struct的两个实例，则所有字段都必须可比较
//	4) 如果是数组，则元素必须可比较。

type A interface{}

type B interface{}

type C interface {
	run()
}

type D interface {
	run()
	do()
}

type CImpl struct{}

func (c CImpl) run() {}

type DImpl struct{}

func (d DImpl) run() {}
func (d DImpl) do()  {}

type Person struct {
	age  int
	name string
}

type Country struct {
	cMap map[string]int
}

func TestComparable(t *testing.T) {
	// comparable ================================================

	fmt.Printf("int %v\n", 3 == 2)
	//fmt.Printf("int8 %v\n", (int8(3) == int8(2)))
	fmt.Printf("int8 %v\n", int8(3) == int8(2))
	fmt.Printf("float64 %v\n", float64(3) == float64(2))
	fmt.Printf("string %v\n", "abc" == "abc")

	fmt.Printf("empty interface %v\n", A("a") == B("a"))
	var c C = CImpl{}
	var d D = DImpl{}
	fmt.Printf("non-empty interface %v\n", c == d)

	fmt.Printf("array %v\n", [2]int{1, 2} == [2]int{1, 2})
	fmt.Printf("named struct %v\n", Person{age: 11, name: "qq"} == Person{age: 11, name: "qq"})
	sn1 := struct {
		age  int
		name string
	}{age: 11, name: "qq"}
	sn2 := struct {
		age  int
		name string
	}{age: 11, name: "qq"}
	fmt.Printf("anonymous struct %v\n", sn1 == sn2) //struct中的所有成员都可以比较的

	//incomparable ================================================
	//IDE直接就报错了
	//fmt.Printf("func %v\n", (func() {} == func() {}))				//func
	//fmt.Printf("slice %v\n", ([]int{1, 2} == []int{1, 2}))		//slice
	//var m1 = make(map[string]int)
	//var m2 = make(map[string]int)
	//fmt.Printf("map %v\n", (m1 == m2))							//map
	//fmt.Printf("struct %v\n", (Country{m1} == Country{m2}))		//struct实例中所有字段都可以则struct实例才能比较
}
