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

const (
	UserPunchRecordPrefix    = "user_punch_record_"
	UserConnectPlatformTable = "user_connect_platform_9"
	NoBreakSign              = (1 << 7) - 1
	UserLoveSpaceTable       = "user_love_space"
	Batch                    = 500
	LogLevel                 = logger.Info
	//LogLevel = logger.Warn
)

type UserPunchRecord struct {
	LoveSpaceId    int64 `gorm:"primaryKey"`
	UserId         int64 `gorm:"primaryKey"`
	BeginDateYear  int32 `gorm:"primaryKey"`
	BeginDateMonth int32 `gorm:"primaryKey"`
	PunchData      int32
	RestoreData    int32
}

type PairCounter struct {
	selfCounter  int32
	loverCounter int32
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

	//date -> selfCounter, loverCounter
	//00-49
	prevHalfGroupCounter := make(map[int32]*PairCounter)
	//50-99
	nextHalfGroupCounter := make(map[int32]*PairCounter)

	var beginnerSid int64
	var userPunchRecords []UserPunchRecord
	var punchRecordTableName string
	var minSpaceIdRecord UserPunchRecord

	for i := 0; i < 100; i++ { //100张分表
		punchRecordTableName = UserPunchRecordPrefix + strconv.Itoa(i)
		log.Printf("====== count table %s", punchRecordTableName)

		dpDb.Table(punchRecordTableName).
			Select("min(love_space_id) as love_space_id").
			Find(&minSpaceIdRecord)
		if minSpaceIdRecord.LoveSpaceId <= 0 {
			log.Println("no punch data in " + punchRecordTableName)
			continue
		}
		beginnerSid = minSpaceIdRecord.LoveSpaceId

		for {
			dpDb.Table(punchRecordTableName).
				Select("love_space_id, user_id, begin_date_year, begin_date_month, punch_data, restore_data").
				Where("love_space_id >= ? and love_space_id < ? and begin_date_year = 2023 and begin_date_month in (3,4)", beginnerSid, beginnerSid+Batch).
				Find(&userPunchRecords)
			log.Println("records len: ", len(userPunchRecords))
			if len(userPunchRecords) == 0 {
				break
			}

			//1 数据过滤，过滤出微信小程序的数据
			uids := make([]int64, len(userPunchRecords))
			for idx, record := range userPunchRecords {
				uids[idx] = record.UserId
			}
			//SELECT user_connect_platform_9.user_id FROM `user_connect_platform_9` left join user_love_space on user_connect_platform_9.user_id = user_love_space.user_id
			//WHERE user_connect_platform_9.user_id in (10001,10001) and user_love_space.status = 0
			var ucpDataArr []UserConnectPlatformData
			wlDb.Table(UserConnectPlatformTable).Select("user_connect_platform_9.user_id").
				Joins("left join user_love_space on user_connect_platform_9.user_id = user_love_space.user_id").
				Where("user_connect_platform_9.user_id in ? and user_love_space.status = ?", uids, 0).
				Find(&ucpDataArr)
			validUidMap := make(map[int64]struct{})
			for _, data := range ucpDataArr {
				validUidMap[data.UserId] = struct{}{}
			}
			//微信小程序的签到数据
			//userId -> month -> UserPunchRecord
			validPunchRecordsMap := make(map[int64]map[int32]UserPunchRecord)
			//spaceId -> pair UserId
			sid2PairUidMap := make(map[int64][]int64)
			for _, record := range userPunchRecords {
				if _, ok := validUidMap[record.UserId]; ok {
					userRecordMap, ok2 := validPunchRecordsMap[record.UserId]
					if ok2 {
						userRecordMap[record.BeginDateMonth] = record
					} else {
						userRecordMap = make(map[int32]UserPunchRecord)
						userRecordMap[record.BeginDateMonth] = record
						validPunchRecordsMap[record.UserId] = userRecordMap
					}
					pairUids, ok3 := sid2PairUidMap[record.LoveSpaceId]
					if ok3 {
						if record.UserId != pairUids[0] {
							pairUids = append(pairUids, record.UserId)
							sid2PairUidMap[record.LoveSpaceId] = pairUids
						}
					} else {
						pairUids = make([]int64, 1)
						pairUids[0] = record.UserId
						sid2PairUidMap[record.LoveSpaceId] = pairUids
					}
				}

				//if record.LoveSpaceId > beginnerSid {
				//	beginnerSid = record.LoveSpaceId
				//}
			}

			//2 数据统计
			for _, userRecordMap := range validPunchRecordsMap {
				//userRecordMap是某个用户两个月的签到数据
				log.Println(userRecordMap)

				//拼接自己的签到数据
				var selfActualSign int32
				var selfSign int32
				var spaceId int64
				var selfUserId int64
				if record, ok := userRecordMap[4]; ok {
					selfActualSign |= record.PunchData << 8
					selfSign |= (record.PunchData | record.RestoreData) << 8
					spaceId = record.LoveSpaceId
					selfUserId = record.UserId
				}
				if record, ok := userRecordMap[3]; ok {
					selfActualSign |= record.PunchData >> 23
					selfSign |= (record.PunchData | record.RestoreData) >> 23
					spaceId = record.LoveSpaceId
					selfUserId = record.UserId
				}
				log.Printf("userId:%v, actualSign: %b\n", selfUserId, selfActualSign)
				log.Printf("selfSign: %b\n", selfSign)

				counter := prevHalfGroupCounter
				if spaceId != 0 && spaceId%100 >= 50 {
					counter = nextHalfGroupCounter
				}

				//拼接另一半的签到数据
				var loverId int64
				var loverSign int32
				for _, userId := range sid2PairUidMap[spaceId] {
					if userId != selfUserId {
						loverId = userId
					}
				}
				if loverId != 0 {
					loverRecordMap := validPunchRecordsMap[loverId]
					log.Println(loverRecordMap)

					if record, ok := loverRecordMap[4]; ok {
						loverSign |= (record.PunchData | record.RestoreData) << 8
					}
					if record, ok := loverRecordMap[3]; ok {
						loverSign |= (record.PunchData | record.RestoreData) >> 23
					}
				}
				log.Printf("loverId:%v, loverSign: %b\n", loverId, loverSign)

				for i := int32(17); i >= 6; i-- {
					if (selfActualSign>>i)&1 == 1 { //当天有签到
						if ((selfSign >> (i - 6)) & NoBreakSign) != NoBreakSign { //自己7天内有断签
							if pairCounter, ok := counter[i]; ok {
								pairCounter.selfCounter += 1
							} else {
								counter[i] = &PairCounter{1, 0}
							}
							log.Printf("%v", counter)
						}
						if ((loverSign >> (i - 6)) & NoBreakSign) != NoBreakSign { //对方7天内有断签
							if pairCounter, ok := counter[i]; ok {
								pairCounter.loverCounter += 1
							} else {
								counter[i] = &PairCounter{0, 1}
							}
							log.Printf("%v", counter)
						}
					}
				}
			}

			beginnerSid += Batch
			log.Println("beginnerSid: ", beginnerSid)
		}

	}

	log.Println("====== Done!")
	log.Println("统计结果：00-49：")
	for idx, counter := range prevHalfGroupCounter {
		log.Printf("%d -> %d, %d\n", idx, counter.selfCounter, counter.loverCounter)
	}
	log.Println("统计结果：50-99：")
	for idx, counter := range nextHalfGroupCounter {
		log.Printf("%d -> %d, %d\n", idx, counter.selfCounter, counter.loverCounter)
	}
}

//2023/04/11 18:54:35 1111101 11111000000
//2023/04/11 18:54:35 1111111 11111100000
//2023/04/11 18:54:35 0000000 00011111110
//
//17 0 1
//16 0 1
//15 0 1
//14 0 1
//13 0 1
//11 0 1
//10 1 1
//9 1 1
//8 1 1
//7 1 0
//6 1 1

//2023/04/11 18:54:35 0000000 00011 111110
//2023/04/11 18:54:35 0000000 00011 111110
//2023/04/11 18:54:35 1111111 11111100000
//
//7 0 1
//6 1 1
