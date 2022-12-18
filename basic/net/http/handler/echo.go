// Package handler
// Echo 隐式实现了 net/http　的 Handler 接口
package handler

import (
	"log"
	"net/http"
)

type Echo struct{}

// 类似类的静态方法
func (h Echo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, "inside echo handler, request detail: ")
	// 请求方式：GET POST DELETE PUT UPDATE
	log.Println("method:", r.Method)
	log.Println("url:", r.URL.Path)
	log.Println("header:", r.Header)
	log.Println("body:", r.Body)
	// 回复
	_, err := w.Write([]byte("Hello " + r.URL.Path[1:]))
	if err != nil {
		log.Println(err)
	}
}
