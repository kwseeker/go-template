# 日志系统

## 基于标准库log实现

log 库中的函数使用的是默认的std对象(标准错误输出：无前缀、输出datetime)

```
var std = New(os.Stderr, "", LstdFlags)
```



## 第三方日志框架

根据 GitHub 热度排名，前三名分别是：

+ logrus
+ zap
+ zerolog

### [logrus](https://github.com/Sirupsen/logrus/blob/master/README.md)

基本使用参考官方README.md就可以。

或者参考 pkg.go.dev [logrus](https://pkg.go.dev/github.com/sirupsen/logrus#section-readme)（额外添加了API文档）, 或者别人的中文博客: [logrus日志使用教程](https://mojotv.cn/2018/12/27/golang-logrus-tutorial)

特性：

+ **是一个可插拔的,结构化的日志框架，全兼容golang标准库日志模块**

+ **拥有六种日志级别**

  Trace, Debug, Info, Warning, Error, Fatal and Panic.

+ **可拓展的Hook机制**

  在初始化时可为logrus添加hook，以实现各种扩展功能。

  原理是每次写入日志时拦截(和很多插拔式拓展组件的原理都一样)，修改logrus.Entry。

  ```
  // logrus在记录Levels()返回的日志级别的消息时会触发HOOK,
  // 按照Fire方法定义的内容修改logrus.Entry.
  type Hook interface {
      Levels() []Level
      Fire(*Entry) error
  }
  ```

  这个功能很有用，比如：

  + 用于报警
    + 将日志额外发到到Email，参考: [logrus_mail](https://github.com/zbindenren/logrus_mail)
    + 将日志额外发到机器人(比如企业微信机器人、钉钉机器人 [dingrus](https://github.com/dandans-dan/dingrus) )
  + 用于日志分割
    + 比如借助 file-rotatelogs
  + ...

+ **可选的日志输出格式**

  logrus内置的格式化器是 logrus.JSONFormatter 和 logrus.TextFormatter。

  分别将日志格式化为JSON & Text。

  另外还支持第三方的格式化器，如 [FluentdFormatter](https://github.com/joonix/log)、[logstash](https://github.com/bshuster-repo/logrus-logstash-hook)等。

+ **Field机制**

  鼓励通过日志字段进行精细化的、结构化的日志记录，而不是冗长的、无法解析的错误消息。

  即下面第一种写法比第二种写法更易于分析。

  ```
  time="2022-03-26T16:31:03+08:00" level=warning msg="The ice breaks!" number=100 omg=true
  time="2022-03-26T16:31:03+08:00" level=warning msg="The ice breaks! number=100 omg=true"
  ```

其他：

+ **Entries**

  就是Field机制存储结构化数据的map，默认有的fields包括：time、msg、level。

+ **Environments**

  logrus默认没有环境的概念。如果某日志在某个特定的环境（比如测试环境）才输出，需要自行实现。

+ **Logger作为io.Writer**

+ **对数旋转**

+ **工具**

  + **[Logrus Mate](https://github.com/gogap/logrus_mate)**

    Logrus mate是Logrus用来管理日志记录器的工具，你可以通过配置文件初始化日志记录器的级别、钩子和格式化器，日志记录器将在不同的环境中以不同的配置生成。

  + **[Logrus Viper Helper](https://github.com/heirko/go-contrib/tree/master/logrusHelper)**

+ **Fatal 处理**
+ **线程安全**