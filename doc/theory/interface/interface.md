# Go interface{} 

interface 是一种类型，没有方法的 interface 类型是 empty interface。

Go 中一个类型实现了一个 interface 中所有方法，就说该类型实现了该 interface, 所以所有类型都实现了 empty interface。

如果一个类型实现了某 interface, 此类型对象就可以通过该 interface 类型传参。

判断 interface 变量存储的是哪种类型，可以通过 `value, ok := em.(T)` 形式做类型断言，em 是 interface 类型的变量，T代表要断言的类型，value 是 interface 变量存储的值，ok 是 bool 类型表示是否为该断言的类型 T。 `em.(T)` 类似 Java 中 `instanceof `关键字。

使用注意事项：

+ **empty interface 可以接受任何类型的参数，但是 `interface{}` 类型的 slice 不可以接受任何类型的 slice**

   `interface{}` 会占用两个字长的存储空间，一个是自身的 methods 数据，一个是指向其存储值的指针，也就是 interface 变量存储的值，因而 slice []interface{} 其长度是固定的`N*2`，但是 []T 的长度是`N*sizeof(T)`，两种 slice 实际存储值的大小是有区别的。 

+ **empty interface 可以接受任何类型的参数，但是不能把某类型参数强转成 empty interface**

  empty interface 和 Java 中的 Object 对象还是不一样。初学 Go 时将 empty interface 当作 Java 的 Object 是错误的。

  ```go
  func TestTypeConvertToInterface(t *testing.T) {
  	arvin := &Person{
  		Name: "Arvin",
  		Age:  18,
  	}
  
  	var obj interface{} = arvin
  	log.Println("obj: ", obj)
  
  	//obj, ok := arvin.(interface{})	//empty interface 可以接受任何类型的参数，但不能把某类型参数强转成 empty interface
  	//语法检查时就会报：Invalid type assertion: arvin.(interface{}) (non-interface type *Person on the left)
  }
  ```

  

