package main

import (
	"fmt"

	"github.com/Vergangenheit/codecrafters-redis-go/app"
	"github.com/Vergangenheit/codecrafters-redis-go/client"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	key      string
	value    string
	hostPort string
)

func init() {
	// Add flags specific to the server command
	setCmd.Flags().StringVar(&key, "key", "", "Key to set")
	setCmd.Flags().StringVar(&value, "value", "", "Value to set")
	setCmd.Flags().StringVar(&hostPort, "hostport", "lovalhost:6379", "host port of the server")

	getCmd.Flags().StringVar(&key, "key", "", "Key to set")
	getCmd.Flags().StringVar(&hostPort, "hostport", "lovalhost:6379", "host port of the server")

	// Bind flags to Viper
	viper.BindPFlag("key", setCmd.Flags().Lookup("key"))
	viper.BindPFlag("value", setCmd.Flags().Lookup("value"))
	viper.BindPFlag("hostport", setCmd.Flags().Lookup("hostport"))

	viper.BindPFlag("key", getCmd.Flags().Lookup("key"))
	viper.BindPFlag("hostport", getCmd.Flags().Lookup("hostport"))
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Send SET command to the server",
	Long:  `send set command to the server with the specified key and value.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := hclog.New(&hclog.LoggerOptions{
			Name:  "client",
			Level: hclog.LevelFromString("INFO"),
		})
		// Retrieve values from Viper
		key := viper.GetString("key")
		value := viper.GetString("value")
		hostPort := viper.GetString("hostport")

		logger.Info(fmt.Sprintf("Sending SET command to server on %s...\n", hostPort))

		// send command
		cl, err := client.NewRedisClient(hostPort)
		if err != nil {
			logger.Error("error instantiating client", err)
			return
		}
		defer cl.Close()

		req := &app.Request{
			Command: app.SET,
			Args:    []string{key, value},
		}
		resp, err := cl.Send(req)
		if err != nil {
			logger.Error("error sending request", err)
			return
		}

		for _, r := range resp {
			fmt.Println(r)
		}

	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Send GET command to the server",
	Long:  `send GET command to the server with the specified key.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := hclog.New(&hclog.LoggerOptions{
			Name:  "client",
			Level: hclog.LevelFromString("INFO"),
		})
		// Retrieve values from Viper
		key := viper.GetString("key")
		hostPort := viper.GetString("hostport")

		logger.Info(fmt.Sprintf("Sending GET command to server on %s...\n", hostPort))

		// send command
		cl, err := client.NewRedisClient(hostPort)
		if err != nil {
			logger.Error("error instantiating client", err)
			return
		}
		defer cl.Close()

		req := &app.Request{
			Command: app.GET,
			Args:    []string{key},
		}
		resp, err := cl.Send(req)
		if err != nil {
			logger.Error("error sending request", err)
			return
		}

		for _, r := range resp {
			fmt.Println(r)
		}

	},
}
