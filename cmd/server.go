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
	port       string
)

func init() {
	// Add flags specific to the server command
	serverStartCmd.Flags().StringVar(&dir, "dir", "/tmp/redis-files", "Directory path for the server")
	serverStartCmd.Flags().StringVar(&dbfilename, "dbfilename", "dump.rdb", "Database filename for the server")
	serverStartCmd.Flags().StringVar(&port, "port", "6379", "port to run server from")

	// Bind flags to Viper
	viper.BindPFlag("dir", serverStartCmd.Flags().Lookup("dir"))
	viper.BindPFlag("dbfilename", serverStartCmd.Flags().Lookup("dbfilename"))
	viper.BindPFlag("port", serverStartCmd.Flags().Lookup("port"))
}

var serverStartCmd = &cobra.Command{
	Use:   "server-start",
	Short: "Start the server",
	Long:  `Start the server with the specified configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve values from Viper
		dir := viper.GetString("dir")
		dbfilename := viper.GetString("dbfilename")
		port := viper.GetString("port")

		fmt.Printf("Starting server on port %s...\n", port)
		fmt.Printf("Using directory: %s\n", dir)
		fmt.Printf("Using database file: %s\n", dbfilename)

		// Add your server startup logic here
		config := &app.Config{
			Dir:        dir,
			DbFilename: dbfilename,
			Port:       port,
		}
		err := app.RunServer(config)
		if err != nil {
			log.Fatal(err)
		}
	},
}
