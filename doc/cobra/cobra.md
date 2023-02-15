# Cobra

参考资料：

+ https://cobra.dev/
+ https://github.com/spf13/cobra

Cobra 是一个用于创建强大的现代 CLI 应用程序的库。 Cobra 用于许多 Go 项目，例如 Kubernetes、Hugo 和 GitHub CLI 等等。

**Cobra功能**：

+ 简单的基于子命令的命令行：app server、app fetch 等等
+ 完全符合POSIX的标志（包含短版本和长版本）
+ 嵌套子命令
+ 全局、本地和级联的标志
+ 使用 cobra init appname和cobra add cmdname 可以很容易生成应用程序和命令
+ 智能提示（app srver... did you mean app server?）
+ 自动生成命令和标志
+ 自动识别 -h --help 等等为help标志
+ 为应用程序自动shell补全（bash、zsh、fish、powershell）
+ 为应用程序自动生成手册
+ 命令别名
+ 灵活定义帮助、用法等等
+ 可选的与viper的紧密集成

**安装与使用**：

```shell
go get -u github.com/spf13/cobra
import "github.com/spf13/cobra"
```

**[Cobra生成器](https://github.com/spf13/cobra-cli/blob/main/README.md)**：

可以帮助用户快速生成基本代码。

```shell
go install github.com/spf13/cobra-cli@latest
cobra-cli init cli-server
# 默认生成的项目结构
cli-server
├── cmd
│   └── root.go
├── LICENSE
└── main.go
```

**Cobra Demo**:

参考：cobra/cli-app。

**主要功能**：

+ 标志（Flags）

  + 持久标志

    所有子命令都会继承的标志。

  + 本地标志

    只应用于当前命令。

    但是可以通过 Command.TraverseChildren 属性，解析父命令的本地标志。

  + 将标志与配置进行绑定

    如：

    ```go
    rootCmd.PersistentFlags().StringVar(&author, "author", "YOUR NAME", "Author name for copyright attribution")
    viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
    ```

    > viper是用于专门解析配置文件的，经常和cobra搭配使用。

  + 标记必传的标志

    ```go
    rootCmd.MarkFlagRequired("region")
    ```

+ 参数校验（Args）

  可以使用命令的Args族字段指定位置参数的验证。

+ 帮助命令

  可以自定义也可以使用默认的。

  ```go
  cmd.SetHelpCommand(cmd *Command)
  cmd.SetHelpFunc(f func(*Command, []string))
  cmd.SetHelpTemplate(s string)
  ```

+ 运行前运行后钩子

  可以在命令的主运行函数 Run() 之前或之后运行函数。PersistentPreRun和PreRun函数将在运行之前执行。PersistentPostRun和PostRun将在运行后执行。如果子函数不声明自己的函数，则它们将继承Persistent*Run函数。

  这些函数按以下顺序运行：

  - PersistentPreRun
  - PreRun
  - Run
  - PostRun
  - PersistentPostRun

+ 未知命令建议

  发生拼写错误时Cobra将自动打印建议，类似git命令。

+ shell补全

