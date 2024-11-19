package main

import (
	"fmt"
	"log"

	"github.com/Vergangenheit/codecrafters-redis-go/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	dir        string
	dbfilename string
)

func init() {
	// Add flags specific to the server command
	serverStartCmd.Flags().StringVar(&dir, "dir", "/tmp/redis-files", "Directory path for the server")
	serverStartCmd.Flags().StringVar(&dbfilename, "dbfilename", "dump.rdb", "Database filename for the server")

	// Bind flags to Viper
	viper.BindPFlag("dir", serverStartCmd.Flags().Lookup("dir"))
	viper.BindPFlag("dbfilename", serverStartCmd.Flags().Lookup("dbfilename"))
}

var serverStartCmd = &cobra.Command{
	Use:   "server-start",
	Short: "Start the server",
	Long:  `Start the server with the specified configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve values from Viper
		dir := viper.GetString("dir")
		dbfilename := viper.GetString("dbfilename")

		fmt.Printf("Starting server...\n")
		fmt.Printf("Using directory: %s\n", dir)
		fmt.Printf("Using database file: %s\n", dbfilename)

		// Add your server startup logic here
		config := &app.Config{
			Dir:        dir,
			DbFilename: dbfilename,
		}
		err := app.RunServer(config)
		if err != nil {
			log.Fatal(err)
		}
	},
}
