package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gar",
	Short: "GAR - Google Agent Runtime CLI",
	Long: `Gar is a CLI tool for managing agent orchestrator sessions.
It provides commands to trigger sessions, resume from checkpoints,
inspect session state, register agents, and run the controller server.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(triggerCmd)
	rootCmd.AddCommand(inspectCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(serveCmd)
}
