package oop

import "testing"

func Test_Inherit(t *testing.T) {
	a := subA{base{}}
	a.run()
	a.sleep()

	b := subB{&base{}}
	b.run()
	b.sleep()

	// 指针实例[字段继承]
	pa := &subA{base{}}
	pa.run() // 3
	pa.sleep()
	// 指针实例[字段指针继承]
	pb := &subB{&base{}}
	pb.run() // 4
	pb.sleep()
}

func Test_Inherit2(t *testing.T) {
	var i iter

	// 如果需要将a实例转化为接口，必须实现接口
	// Base结构体已经实现了接收者为实例的sleep()方法
	// 那么可以在subA结构体实现 接收者为实例接受的run()方法即可
	// a实例[字段继承]
	//a := subA{base{}}
	//i = a   // A1	error subA未实现接口，父类仅仅实现了sleep()方法,run()没有实现
	//i.run() // A1
	//i.sleep()

	// b实例[字段指针继承]
	b := subB{&base{}}
	i = b   // subB实现了接口
	i.run() // B1
	i.sleep()

	// =======指针实例转换为接口==========
	// 指针实例[字段继承]
	pa := &subA{base{}}
	i = pa  // subA实现了接口
	i.run() // C1
	i.sleep()

	// 指针实例[字段指针继承]
	pb := &subB{&base{}}
	i = pb  //  subB实现接口
	i.run() //  D1
	i.sleep()
}
