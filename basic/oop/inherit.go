package oop

import "fmt"

//值类型调用 指针类型调用
//值接收者 指针接收者

type iter interface {
	run()
	sleep()
}

type base struct{}

// p 是指针接收者
func (p *base) run() {
	fmt.Println("Base::run()")
}

// p 是值接收者
func (p base) sleep() {
	fmt.Println("Base::sleep()")
}

type subA struct {
	base //值继承
}

type subB struct {
	*base //指针继承
}
