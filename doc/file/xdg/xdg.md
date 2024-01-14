# xdg

xdg 是一个规范([XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html))，该规范定义了一组用于存储应用程序文件(包括数据和配置文件)的标准路径。

出于可移植性和灵活性的考虑，应用程序应该使用XDG定义的位置，而不是硬编码路径。该包还包括知名用户目录的位置，以及其他常见目录，如字体和应用程序。

此处介绍 `xdg` 规范的 Go 实现：https://github.com/adrg/xdg

[xdg doc](https://pkg.go.dev/github.com/adrg/xdg)