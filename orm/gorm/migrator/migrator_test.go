package migrator

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"testing"
)

/*
GORM 迁移测试

比如有下面一个表，想要加入一个 gender 字段

CREATE TABLE `test`.`user` (
  `id` INT(11) NOT NULL,
  `name` VARCHAR(45) NULL,
  `age` TINYINT(4) NULL,
  PRIMARY KEY (`id`));
*/

type User struct {
	ID     int32
	Name   string `gorm:"default:someone"`
	Age    int8   `gorm:"default:18"`
	Gender int8   `gorm:"default:0"`
}

func TestAutoMigrate(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("migrate failed, err=%v", err)
		return
	}
	log.Println("migrate done")
}
