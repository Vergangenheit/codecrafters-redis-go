package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverStartCmd)
}

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "A brief description of your CLI",
	Long:  `A longer description for your CLI tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to MyCLI! Use 'server start' to start the server.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
