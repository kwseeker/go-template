package main

import (
	"fmt"
	"github.com/ip2location/ip2location-go/v9"
)

func main() {
	//二进制数据库
	db, err := ip2location.OpenDB("/home/lee/tool/IPV6-COUNTRY-REGION-CITY.BIN")
	if err != nil {
		fmt.Print(err)
		return
	}
	//支持IPv6
	ip := "8.8.8.8"
	results, err := db.Get_all(ip)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("result: %v\n", results)

	db.Close()
}
