package reflect

type User struct {
	Id   int    `validate:"number,min=1,max=1000"`
	Name string `validate:"string,min=2,max=10"`
	//Bio   string `validate:"string"`
	Email string `validate:"email"`
}

//type User struct {
//	Id    int    `validate:"-"`
//	Name  string `validate:"presence,min=2,max=32"`
//	Email string `validate:"email,required"`
//}
