package main

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"strconv"
	"time"
)

/*
统计补签卡数大于0的用户数量（过滤掉其他平台的用户），将用户信息打印到文件
*/

const (
	UserRestoreCardTablePrefix = "user_restore_card_"
	UserConnectPlatformTable   = "user_connect_platform_9"
	UserLoveSpaceTable         = "user_love_space"
	Batch                      = 500
	//LogLevel                   = logger.Info
	LogLevel = logger.Warn
)

type UserRestoreCard struct {
	LoveSpaceId int64 `gorm:"primaryKey"`
	UserId      int64 `gorm:"primaryKey"`
	RestoreCard int32
}

type UserConnectPlatformData struct {
	UserId int64 `gorm:"primaryKey"`
}

func createConn(dsn string) (*gorm.DB, *sql.DB) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(LogLevel),
	})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(4)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, sqlDB
}

func main() {
	//DEV
	dpDsn := "root:123456@tcp(127.0.0.1:3306)/daily_punch?charset=utf8mb4&parseTime=True&loc=Local"
	wlDsn := "root:123456@tcp(127.0.0.1:3306)/welove?charset=utf8mb4&parseTime=True&loc=Local"

	dpDb, _ := createConn(dpDsn)
	wlDb, _ := createConn(wlDsn)

	counter := make(map[int32]int32)

	var beginnerSid int64
	var restoreCards []UserRestoreCard
	var restoreCardTableName string
	for i := 0; i < 100; i++ { //100张分表
		beginnerSid = 0
		restoreCardTableName = UserRestoreCardTablePrefix + strconv.Itoa(i)
		log.Printf("====== count table %s", restoreCardTableName)

		for {
			dpDb.Table(restoreCardTableName).
				Select("love_space_id, user_id, restore_card").
				Where("love_space_id > ? and restore_card > ?", beginnerSid, 0).
				Limit(Batch).Find(&restoreCards)

			//数据过滤
			uids := make([]int64, len(restoreCards))
			for idx, card := range restoreCards {
				uids[idx] = card.UserId
			}
			//select ucp.user_id from user_connect_platform_9 ucp left join user_love_space uls on ucp.user_id = uls.user_id
			//where ucp.user_id in (10001, 10002, 1003) and uls.status = 0;
			var ucpDataArr []UserConnectPlatformData
			wlDb.Table(UserConnectPlatformTable).Select("user_connect_platform_9.user_id").
				Joins("left join user_love_space on user_connect_platform_9.user_id = user_love_space.user_id").
				Where("user_connect_platform_9.user_id in ? and user_love_space.status = ?", uids, 0).
				Find(&ucpDataArr)
			validUidMap := make(map[int64]struct{})
			for _, data := range ucpDataArr {
				validUidMap[data.UserId] = struct{}{}
			}
			validRestoreCards := make([]UserRestoreCard, len(ucpDataArr))
			i2 := 0
			for _, card := range restoreCards {
				if _, ok := validUidMap[card.UserId]; ok {
					validRestoreCards[i2] = card
					i2++
				}

				if card.LoveSpaceId > beginnerSid {
					beginnerSid = card.LoveSpaceId
				}
			}

			//统计
			for _, card := range validRestoreCards {
				log.Println(card)

				if card.RestoreCard >= 100 {
					counter[100]++
				} else {
					counter[card.RestoreCard]++
				}
			}

			log.Println("beginnerSid: ", beginnerSid)

			if len(restoreCards) < Batch {
				break
			}
		}

	}

	log.Println("====== Done!")
	log.Println("统计结果：", counter)
}
