// Copyright © 2017 Philippe Hässig <phil@neckhair.ch>

package cmd

import (
	"fmt"
	"github.com/neckhair/owntracks-eventr/listener"
	"github.com/spf13/cobra"
)

var url string
var filename string

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for events and write them into a file",
	Long:  `Listen for events and write them into a file line by line.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("--> Listening on %s for events, write to %s\n\n", url, filename)
		listener.Listen(url)
	},
}

func init() {
	RootCmd.AddCommand(listenCmd)

	listenCmd.Flags().StringVarP(&url, "url", "u", "tcp://localhost:1883", "Connection string")
	listenCmd.Flags().StringVarP(&filename, "output", "o", "eventlog.txt", "Path to destination file")
}
