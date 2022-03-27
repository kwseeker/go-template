package log

import (
	log "github.com/sirupsen/logrus"
	"io"
	file2 "kwseeker.top/kwseeker/go-template/file"
	"os"
)

type DefaultFieldHook struct {
}

// 这时实例化主类的方式，这里不适合
//func New() *DefaultFieldHook {
//	return &DefaultFieldHook{}
//}
// 实例化内部子类的方式
func newDefaultFiledHook() *DefaultFieldHook {
	return &DefaultFieldHook{}
}

func (hook *DefaultFieldHook) Fire(entry *log.Entry) error {
	entry.Data["app_name"] = "logrusSample"
	return nil
}

func (hook *DefaultFieldHook) Levels() []log.Level {
	return log.AllLevels
}

// 自定义Logger
var formatLogger = log.New()
var loggerWithMethod = log.New()
var loggerCustomized = log.New()
var loggerWithHook = log.New()

func init() {
	loggerWithMethod.SetReportCaller(true)

	path := "/tmp/go/test/logrus.log"
	file2.CreateFileWithParentDir(path)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatalln("Failed to open logger file:", err)
		return
	}
	//loggerCustomized.Out = file
	loggerCustomized.SetOutput(io.MultiWriter(file, os.Stdout))
	loggerCustomized.SetLevel(log.WarnLevel)

	//loggerWithHook.AddHook(&DefaultFieldHook{})
	loggerWithHook.AddHook(newDefaultFiledHook())
}

// Print 基本使用
func Print() {
	log.Info("A walrus appears")

	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")

	log.WithFields(log.Fields{
		"omg":    "true",
		"number": 100,
	}).Warn("The ice breaks!")

	loggerWithMethod.Info("A walrus appears")

	loggerCustomized.Warn("A walrus appears")

	//logger with default fields
	requestId := 10001
	userIp := "192.168.1.101"
	loggerWithDefaultField := log.WithFields(log.Fields{"request_id": requestId, "user_ip": userIp})
	loggerWithDefaultField.Info("something happened on that request")
	loggerWithDefaultField.Warn("something not great happened")

	loggerWithHook.Info("log by hook")
}

// FormatPrint 格式化输出
func FormatPrint() {
	// 默认是Text格式
	formatLogger.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")

	//改为JSON格式
	formatLogger.SetFormatter(&log.JSONFormatter{})
	formatLogger.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
}
