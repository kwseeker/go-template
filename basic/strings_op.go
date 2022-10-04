package main

import (
	"log"
	"strings"
)

func main() {
	rest := "//localhost:9003/echo?param=100"
	var rawQuery string
	rest, rawQuery, _ = strings.Cut(rest, "?")
	log.Print(rest)
	log.Print(rawQuery)
}
