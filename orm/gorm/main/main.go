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

const TableName = "ring_box"
const Batch = 10000

type void struct{}

type RingBox struct {
	Id          uint `gorm:"primary"`
	UserId      uint64
	ItemId      uint
	State       uint8
	Display     uint8
	BindLoverId uint64
	BindRoomId  uint64
	BindTime    *time.Time
}

// TableName 自定义表名，默认是ring_boxes
func (RingBox) TableName() string {
	return TableName
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

	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
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

	var rings []RingBox
	var minUserId uint64
	var maxUserId uint64
	//查最小和最大的userId
	row := db.Table(TableName).Select("min(user_id)").Row()
	err = row.Scan(&minUserId)
	if err != nil {
		return
	}
	row = db.Table(TableName).Select("max(user_id)").Row()
	err = row.Scan(&maxUserId)
	if err != nil {
		return
	}
	log.Printf("min userId: %d, max userId: %d\n", minUserId, maxUserId)

	//batch query
	ptr := minUserId

	exceptSet := make(map[string]void)
	for ptr <= maxUserId {
		log.Println("handle: user_id: ", ptr, ptr+uint64(Batch))
		db.Where("user_id >= ? and user_id < ?", ptr, ptr+uint64(Batch)).Find(&rings)
		ptr += uint64(Batch)
		if len(rings) == 0 {
			continue
		}
		//filter: user_id & bind_lover_id & item_id all same record
		set := make(map[string]void)
		var member void
		for _, ring := range rings {
			key := strconv.FormatUint(ring.UserId, 10) + ":" + strconv.Itoa(int(ring.ItemId)) + ":" + strconv.FormatUint(ring.BindLoverId, 10)
			if _, exists := set[key]; !exists {
				set[key] = member
			} else if _, exists := exceptSet[key]; !exists {
				exceptSet[key] = member
				//repair
				log.Println(">>>>>>>", key)
				var loverRing RingBox
				if err := db.First(&loverRing, "user_id = ? and bind_lover_id = ? and item_id = ?", ring.BindLoverId, ring.UserId, ring.ItemId).Error; err != nil {
					log.Println(">>>>>>> repairing ...")
					db.Model(&ring).Debug().Updates(RingBox{UserId: ring.BindLoverId, BindLoverId: ring.UserId})
				} else {
					if ring.State == 1 {
						log.Println(">>>>>>> lover has this item, not use, delete")
						db.Debug().Delete(&RingBox{}, ring.Id)
					} else {
						log.Println(">>>>>>> lover has this item, on use, skip")
					}
				}
			} else {
				log.Println("+++++++ ", key)
				if ring.State == 1 {
					log.Println("+++++++ delete")
					db.Debug().Delete(&RingBox{}, ring.Id)
				}
			}
		}

		for k := range exceptSet {
			delete(exceptSet, k)
		}
	}

	log.Println("Done!")
}
