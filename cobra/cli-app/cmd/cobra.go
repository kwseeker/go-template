package cmd

import (
	"errors"
	"fmt"
	"kwseeker.top/kwseeker/go-template/cobra/cli-app/cmd/config"
	"kwseeker.top/kwseeker/go-template/cobra/cli-app/cmd/server"
	"kwseeker.top/kwseeker/go-template/cobra/cli-app/cmd/version"
	"kwseeker.top/kwseeker/go-template/cobra/cli-app/common"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	//命令行指令简单使用方法（一行）
	Use: "cli-app",
	//命令别名
	Aliases: []string{"ca"},
	//APP简短描述
	Short: "cli-app",
	//APP详细描述
	Long: `A cli app demo by cobra!`,
	//隐藏默认的使用说明
	SilenceUsage: true,
	Version:      "0.0.1",
	//用于传参校验（如下：只是传参不为空校验）
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New("requires at least one arg")
		}
		return nil
	},
	//
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	//命令业务逻辑
	//Run不带错误返回值，如果产生错误需要返回可用RunE
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
	//RunE: func(cmd *cobra.Command, args []string) error {
	//	if err := someFunc(); err != nil {
	//		return err
	//	}
	//	return nil
	//},
}

func init() {
	//自定义子命令
	rootCmd.AddCommand(version.StartCmd) //其实没必要，cobra.Command 的 Version 参数可以实现这个功能
	rootCmd.AddCommand(config.StartCmd)
	rootCmd.AddCommand(server.StartCmd)
	//自定义帮助命令
	//rootCmd.SetHelpCommand(cmd *Command)
	//rootCmd.SetHelpFunc(f func(*Command, []string))
	//rootCmd.SetHelpTemplate(s string)
}

func tip() {
	usageStr := `cli-app:` + common.Version + `
usage: cli-server 
more details: ......`
	fmt.Printf("%s\n", usageStr)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
