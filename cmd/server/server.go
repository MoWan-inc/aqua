package server

import (
	"context"
	"github.com/MoWan-inc/aqua/pkg/config"
	"github.com/spf13/cobra"
	"os/signal"
	"syscall"
)

// TODO: 添加zlog及其配置，添加日志
func NewCmd() *cobra.Command {
	cfg := config.DefaultServerConfig()

	cmd := &cobra.Command{
		Use:   "server",
		Short: "run mowan server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.Validate(); err != nil {
				return err
			}
			ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer cancel()
			return mainFunc(ctx, cfg)
		},
	}
	cmd.Flags().VarP(cfg, "config-path", "", "config file path")
	return cmd
}

func mainFunc(ctx context.Context, cfg *config.ServerConfig) error {
	return nil
}
