package main

import (
	"os"

	"github.com/scythe504/webtorrent/internal/commands"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
)

var rootCmd = &cobra.Command{
	Use:     "fluxstream",
	Short:   "Fluxstream - Torrent media streamer",
	Long:    `fluxstream-cli is a tool for running the fluxstream server and web on your desktop`,
	Version: version,
}
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Fluxstream server",
	Long:  `Start the Fluxstream server using docker-compose with neat logs and network info`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.Start()
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Sets up required configs and docker engine (if not installed)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.Setup()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(setupCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
