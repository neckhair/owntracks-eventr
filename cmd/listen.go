// Copyright © 2017 Philippe Hässig <phil@neckhair.ch>

package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/neckhair/owntracks-eventr/listener"
)

var config = listener.Configuration{}

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for events and write them into a file",
	Long:  `Listen for events and write them into a file line by line.`,

	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("--> Listening for MQTT events\n")
		fmt.Printf("Server:  %s\n", config.Url)
		fmt.Printf("Output:  %s\n", config.Filename)
		fmt.Printf("Logfile: %s\n\n", viper.GetString("LogFile"))
	},

	Run: func(cmd *cobra.Command, args []string) {
		listener := listener.NewListener(&config)
		listener.TLSConfig.InsecureSkipVerify = viper.GetBool("insecure")

		if err := listener.Start(); err != nil {
			fmt.Println("Could not connect to MQTT server.")
			log.Fatalln(err)
		}
		defer listener.Stop()

		for {
		}
	},
}

func init() {
	RootCmd.AddCommand(listenCmd)

	listenCmd.Flags().StringVarP(&config.Url, "url", "u", "tls://localhost:8883", "Connection string")
	listenCmd.Flags().StringVarP(&config.Filename, "output", "o", "eventlog.txt", "Path to destination file")

	listenCmd.Flags().Bool("insecure", false, "Skip TLS certificate verification")
	viper.BindPFlag("insecure", listenCmd.Flags().Lookup("insecure"))
}
