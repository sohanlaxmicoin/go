package main

import (
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run orbit server",
	Long:  "serve initializes then starts the orbit HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		initApp(cmd, args)
		app.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
