package clicmd

import (
	"log/slog"
	"os"

	"links-sync-go/internal/config"
	webapi "links-sync-go/internal/web_api"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "running server",
	Run: func(cmd *cobra.Command, args []string) {
		// setup default logger
		logLevel := &slog.LevelVar{} // info
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		}))
		slog.SetDefault(logger)

		config, err := config.ReadConfig(cfgFile)
		if err != nil {
			slog.Error("failed to read config",
				slog.Any("error", err),
			)
			os.Exit(1)
		}

		// set log level from config
		err = logLevel.UnmarshalText([]byte(config.Core.LogLevel))
		if err != nil {
			slog.Error("failed to parse log level",
				slog.Any("error", err),
			)
			os.Exit(1)
		}
		server := webapi.NewApiServer(config)

		server.Run()
	},
}
