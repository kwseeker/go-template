package gorm

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

// User
/*
CREATE TABLE `test`.`user` (
`id` INT NOT NULL AUTO_INCREMENT,
`name` VARCHAR(45) NOT NULL,
`email` VARCHAR(45) NULL,
`age` SMALLINT(3) NOT NULL,
`birthday` DATE NOT NULL,
`member_number` VARCHAR(45) NULL,
`activated_at` DATETIME NULL,
`created_at` DATETIME NOT NULL,
`updated_at` DATETIME NOT NULL,
`deleted_at` DATETIME NULL,
PRIMARY KEY (`id`));
*/
type User struct {
	gorm.Model
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
}

// TableName 自定义表名，默认是users
func (User) TableName() string {
	return "user"
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Create
	email := "xiaohuileee@gmail.com"
	birthday, _ := time.Parse("2006-01-02", "2000-01-01")
	current := time.Now()
	userPtr := &User{Model: gorm.Model{CreatedAt: current, UpdatedAt: current}, Name: "Arvin", Email: &email, Age: 18, Birthday: &birthday}
	db.Create(userPtr)
	if db.Error != nil {
		log.Fatalln("Create failed: " + db.Error.Error())
		return
	}
	log.Println("userPtr.ID = ", userPtr.ID)

	// Read
	var user User
	db.First(&user, userPtr.ID)
	log.Println(user)
	db.Find(&user, "name", "Arvin")
	log.Println(user)
	db.First(&user, "name = ?", "Arvin")
	log.Println(user)

	// Update
	db.Model(&user).Update("Age", 20)
	// Update - 更新多个字段
	db.Model(&user).Updates(User{Age: 21, Name: "Arvin Lee"}) // 仅更新非零值字段
	db.Model(&user).Updates(map[string]interface{}{"Age": 22, "Name": "Kwseeker"})

	// Delete - 删除
	db.Delete(&user, user.ID)
}
