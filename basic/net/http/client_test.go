package http

import "testing"

func TestGet(t *testing.T) {
	Get("http://127.0.0.1:8080/go")
	Get("http://127.0.0.1:8080/hello/Arvin")
}

func TestPost(t *testing.T) {
	jsonMap := map[string]interface{}{
		"name": "Arvin",
		"age":  18,
	}
	Post("http://127.0.0.1:8080/json/somemsg", jsonMap)
}
