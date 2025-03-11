package main

import (
	"github.com/MoWan-inc/aqua/cmd/server"
	"github.com/MoWan-inc/aqua/cmd/util"
	"github.com/MoWan-inc/aqua/pkg/config"
	"github.com/MoWan-inc/aqua/pkg/version"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Version: version.String(),
	}
)

func mainFunc() error {
	rootCmd.PersistentFlags().StringVar(&config.LogConfigPath, "log-config-path", "", "log config file path, if empty use debug mode")

	// TODO: 初始化日志相关配置

	rootCmd.AddCommand(server.NewCmd())

	return rootCmd.Execute()
}

func main() {
	err := mainFunc()
	util.ExitPrintFatalError(err)
}
