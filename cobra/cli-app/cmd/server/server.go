package server

import (
	"github.com/spf13/cobra"
	"log"
)

var (
	appName  string
	StartCmd = &cobra.Command{
		Use:     "server",
		Short:   "Create a new server",
		Long:    "Use when you need to create a new server",
		Example: "cli-app server -n admin",
		//Run()之前执行
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		//server命令的业务逻辑
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&appName, "name", "n", "",
		"Start server with provided configuration file")
}

func setup() {
	log.Printf("exec setup before run ...\n")
}

func run() {
	log.Printf("start server %s ...\n", appName)
}
