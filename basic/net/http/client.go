package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Get(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	// 200 OK
	log.Printf("response: Status: %s, Header: %s\n", resp.Status, resp.Header)

	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			log.Println(err)
			return
		}
		log.Println("读取完毕")
		res := string(buf[:n])
		log.Println(res)
		break
	}
}

func Post(url string, body map[string]interface{}) {
	bytesData, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
	}
	reader := bytes.NewReader(bytesData)
	resp, err := http.Post(url, "application/json;charset=UTF-8", reader)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	// 200 OK
	log.Printf("response: Status: %s, Header: %s\n", resp.Status, resp.Header)

	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			log.Println(err)
			return
		}
		log.Println("读取完毕")
		res := string(buf[:n])
		log.Println(res)
		break
	}
}

func Post2(url string, jsonParam []byte) {
	//发送请求
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonParam)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Authorization", "6c7b0b26-3773-59d5-b188-25d22f618f33")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	//响应
	response, err := ioutil.ReadAll(resp.Body)
	fmt.Println("response:", string(response))
	if err != nil {
		log.Println("Read failed:", err)
		return
	}
	log.Println("response from console-approve.success:", string(response))
}

func requestUrl(url string, method string, body io.Reader) ([]byte, error) {
	client := &http.Client{}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=utf-8")
	//处理返回结果
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	buf, err := ioutil.ReadAll(response.Body)
	if response.StatusCode == http.StatusOK {
		return buf, nil
	} else {
		return nil, errors.New(fmt.Sprint("StatusCode=", response.StatusCode, " msg=", string(buf)))
	}
}
