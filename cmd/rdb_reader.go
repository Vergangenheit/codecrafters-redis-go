package main

import (
	"fmt"
	"log"

	"github.com/Vergangenheit/codecrafters-redis-go/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	// Add flags specific to the server command
	readRdbCmd.Flags().StringVar(&dbfilename, "dbfilename", "dump.rdb", "Database filename for the server")

	// Bind flags to Viper
	viper.BindPFlag("dbfilename", serverStartCmd.Flags().Lookup("dbfilename"))
}

var readRdbCmd = &cobra.Command{
	Use:   "read-file",
	Short: "reads and decodes a dump file",
	Long:  `Read and decodes a rdb file and parses into the inmemory store`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve values from Viper
		dbfilename := viper.GetString("dbfilename")

		fmt.Printf("Using database file: %s\n", dbfilename)

		inMem, err := app.ReadRedisDBFile(dbfilename)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(inMem)
	},
}
