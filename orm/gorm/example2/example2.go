package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"time"
)

const UserTableNamePrefix = "user_"
const CoinTableNamePrefix = "house_coin_account_"
const Batch = 100

type void struct{}

type User struct {
	LoveSpaceId  int64     `gorm:"primaryKey"`
	RegisterTime time.Time `gorm:"type:datetime;column:register_time"`
}

type HouseAccount struct {
	LoveSpaceId int64 `gorm:"primaryKey"`
	Balance     int
}

func main() {
	args := os.Args
	var logLevel logger.LogLevel
	if len(args) < 2 {
		log.Println("no params failed! use Warn")
		logLevel = logger.Warn
	} else {
		logLevelParam, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalln("input params failed! ", err.Error())
			return
		}
		logLevel = logger.LogLevel(logLevelParam)
	}

	dsn := "root:123456@tcp(127.0.0.1:3306)/house?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
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

	for i := 0; i < 20; i++ {
		var users []User
		userTableName := UserTableNamePrefix + strconv.Itoa(i)

		for minSpaceId := int64(0); ; {
			//select love_space_id, register_time from user_x where register_time < '' order by love_space_id limit 100;
			db.Table(userTableName).Select("distinct(love_space_id), register_time").
				Where("love_space_id > ? and register_time < ?", minSpaceId, "2022-12-06 17:30:00").
				Limit(Batch).Find(&users)

			if len(users) == 0 {
				break
			}

			//处理
			for _, user := range users {
				var account HouseAccount

				coinTableName := CoinTableNamePrefix + strconv.Itoa(int(user.LoveSpaceId%100))
				db.Table(coinTableName).Select("love_space_id, balance").
					Where("love_space_id = ?", user.LoveSpaceId).Take(&account)
				log.Printf("%v", user)

				if account.LoveSpaceId == 0 && account.Balance == 0 { //不存在则创建
					account = HouseAccount{LoveSpaceId: user.LoveSpaceId, Balance: 5000}
					if err := db.Table(coinTableName).Create(&account).Error; err != nil {
						log.Printf("insert account failed: %v", account)
					} else {
						log.Printf(">>>>>>> insert account success: %v", account)
					}
				} else { //更新
					db.Table(coinTableName).Where("love_space_id", user.LoveSpaceId).Update("balance", account.Balance+5000)
					log.Printf(">>>>>>> update account success + 5000: %v", account)
				}
			}

			user := users[len(users)-1]
			minSpaceId = user.LoveSpaceId
			log.Printf(">>>>>>> next minSpaceId: %d", minSpaceId)
		}
	}

	log.Println("Done!")
}
