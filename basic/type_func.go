package main

import "log"

//自定义函数类型
//语法 type typeName func([paramType]) [returnType]
//拥有相同的形参类型和返回值类型的函数是同一函数类型
//作用：类似Java语言的对象向上造型（父类型引用指向子类型对象），C语言函数指针，可以实现多态。
//应用：参考标准库websocket/server.go

// 下面这两个方法属于同一函数类型
func add(a, b int) int {
	return a + b
}

func subtract(a, b int) int {
	return a - b
}

// CalculateOperator 为上面两个方法定义的一个函数类型
type CalculateOperator func(int, int) int

// 使用一：函数类型的方法（面向对象）
func (op CalculateOperator) doCalculate(a, b int) int {
	return op(a, b)
}

// Calculate 使用二：函数（面向过程）
func Calculate(op CalculateOperator, a int, b int) int {
	return op(a, b)
}

func main() {
	func1 := CalculateOperator(add)
	func2 := CalculateOperator(subtract)

	a := 10
	b := 2

	//面向对象调用
	ret1 := func1.doCalculate(a, b)
	ret2 := func2.doCalculate(a, b)

	//面向过程调用
	ret3 := Calculate(func1, a, b)
	ret4 := Calculate(func2, a, b)

	log.Print(ret1, ret2, ret3, ret4)
}
