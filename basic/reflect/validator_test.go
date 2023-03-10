package reflect

import (
	"fmt"
	"testing"
)

func TestValidate(t *testing.T) {
	//错误格式
	user := User{
		Id:    0,
		Name:  "super long string",
		Email: "foobar",
	}

	//正确格式
	//user := User{
	//	Id:    1,
	//	Name:  "Arvin Lee",
	//	Email: "xiaohuileee@gmail.com",
	//}

	fmt.Println("Errors: ")
	for i, err := range validateStruct(user) {
		fmt.Printf("\t%d. %s\n", i+1, err.Error())
	}
}
