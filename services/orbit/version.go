package main

import (
	"fmt"

	"github.com/spf13/cobra"
	apkg "github.com/rover/go/support/app"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print orbit version",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(apkg.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
