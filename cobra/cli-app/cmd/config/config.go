package config

import (
	"github.com/spf13/cobra"
	"log"
)

// config -c conf/settings.yml

var (
	configYml string //配置文件路径
	StartCmd  = &cobra.Command{
		Use:     "config",
		Short:   "Get Application config info",
		Example: "cli-app config -c conf/settings.yml",
		//config命令的业务逻辑
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml",
		"Start server with provided configuration file")
}

func run() {
	log.Printf("config file path: %s\n", configYml)
}
