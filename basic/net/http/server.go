package http

import (
	"encoding/json"
	"io/ioutil"
	handler2 "kwseeker.top/kwseeker/go-template/basic/net/http/handler"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, req *http.Request) {
	log.Println(req.RemoteAddr, "inside root handler, request detail: ")
	// 请求方式：GET POST DELETE PUT UPDATE
	log.Println("method:", req.Method)
	log.Println("url:", req.URL.Path)
	log.Println("header:", req.Header)
	log.Println("body:", req.Body)
	// 回复
	_, err := w.Write([]byte("Hello " + req.URL.Path[1:]))
	if err != nil {
		log.Println(err)
	}
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte("Hello " + req.URL.Path[len("/hello/"):]))
	if err != nil {
		log.Println(err)
	}
}

type TestStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// json 请求处理
func jsonHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))

	var t TestStruct
	err = json.Unmarshal(body, &t)
	//decoder := json.NewDecoder(req.Body)
	//err = decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t)

	//再原封不动写回去
	_, respErr := w.Write(body)
	if respErr != nil {
		log.Println(respErr)
	}
	//decoder := json.NewDecoder(req.Body)
	//var t test_struct
	//err := decoder.Decode(&t)
	//if err != nil {
	//	panic(err)
	//}
	//defer req.Body.Close()
	//log.Println(t.Test)

	//req.ParseForm()
	//log.Println(req.Form)
	////LOG: map[{"test": "that"}:[]]
	//var t test_struct
	//for key, _ := range req.Form {
	//	log.Println(key)
	//	//LOG: {"test": "that"}
	//	err := json.Unmarshal([]byte(key), &t)
	//	if err != nil {
	//		log.Println(err.Error())
	//	}
	//}
	//log.Println(t.Test)
}

func start() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/hello/", helloHandler)
	http.HandleFunc("/json/", jsonHandler)

	log.Println("server start")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func startEcho() {
	//var echoHandler handler2.Echo

	log.Println("echo server start")
	//err := http.ListenAndServe("localhost:8080", echoHandler)
	err := http.ListenAndServe("localhost:8080", handler2.Echo{})
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
