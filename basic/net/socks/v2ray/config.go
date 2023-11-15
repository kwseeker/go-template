package main

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/infra/conf"
	infraJson "kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/infra/conf/json"
	"kwseeker.top/kwseeker/go-template/basic/os/file"
	"os"
)

// LoadConfig 从 configFile 中读取数据解析到 Config 对象
func LoadConfig(configFile string) (*conf.Config, error) {
	if !file.IsExist(configFile) {
		return nil, errors.New("config file is not exist")
	}

	jsonConfig := &conf.Config{}
	fileBytes, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	bsr := bytes.NewReader(fileBytes)
	contentBuffer := bytes.NewBuffer(make([]byte, 0, 10240))
	// TeeReader() 返回 teeReader 当读取 teeReader 中的内容时，会无缓冲的将读取内容写入到 Writer 中
	jsonReader := io.TeeReader(&infraJson.Reader{
		Reader: bsr,
	}, contentBuffer)
	jsonDecoder := json.NewDecoder(jsonReader)
	if err = jsonDecoder.Decode(jsonConfig); err != nil {
		log.Error("failed to load config file")
		return nil, err
	}
	return jsonConfig, nil
}
