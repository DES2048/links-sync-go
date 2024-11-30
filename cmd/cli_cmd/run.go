package clicmd

import (
	"links-sync-go/internal/config"
	webapi "links-sync-go/internal/web_api"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "runnig server",
	Run: func(cmd *cobra.Command, args []string) {

		config, err := config.ReadConfig(cfgFile)

		if err != nil {
			log.Fatalf("Failed to read config: %s", err)
		}

		server := webapi.NewApiServer(config)

		server.Run()
	},
}
