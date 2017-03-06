// Copyright © 2017 Philippe Hässig <phil@neckhair.ch>

package cmd

import (
	"fmt"
	"github.com/neckhair/owntracks-eventr/listener"
	"github.com/spf13/cobra"
)

type listenCfg struct {
	Url      string
	Filename string
}

var config = listenCfg{}

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for events and write them into a file",
	Long:  `Listen for events and write them into a file line by line.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("--> Listening for MQTT events\n")
		fmt.Printf("Server:\t%s\n", config.Url)
		fmt.Printf("Output:\t%s\n\n", config.Filename)

		listener.Listen(config.Url)
	},
}

func init() {
	RootCmd.AddCommand(listenCmd)

	listenCmd.Flags().StringVarP(&config.Url, "url", "u", "tcp://localhost:1883", "Connection string")
	listenCmd.Flags().StringVarP(&config.Filename, "output", "o", "eventlog.txt", "Path to destination file")
}
