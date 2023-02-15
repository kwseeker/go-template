package version

import (
	"fmt"
	"github.com/spf13/cobra"
	"kwseeker.top/kwseeker/go-template/cobra/cli-app/common"
)

var (
	StartCmd = &cobra.Command{
		Use:     "version",
		Short:   "Get version info",
		Example: "cli-app version",
		PreRun: func(cmd *cobra.Command, args []string) {
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func run() error {
	fmt.Println(common.Version)
	return nil
}
