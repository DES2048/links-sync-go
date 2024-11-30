package clicmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yml", "config file path")
}

var rootCmd = &cobra.Command{
	Use:   "links-sync-server",
	Short: "Links sync backend server",
	Run: func(cmd *cobra.Command, args []string) {

		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
