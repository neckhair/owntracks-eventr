// Copyright © 2017 Philippe Hässig <phil@neckhair.ch>

package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
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

		var err error
		if listener.TLSConfig, err = tlsConfig(); err != nil {
			log.Fatalln(err)
		}

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

	listenCmd.Flags().StringVarP(&config.Url, "url", "", "tls://localhost:8883", "Connection string")
	listenCmd.Flags().StringVarP(&config.Filename, "output", "o", "eventlog.txt", "Path to destination file")
	listenCmd.Flags().StringVarP(&config.Username, "username", "u", "", "MQTT Username")
	listenCmd.Flags().StringVarP(&config.Password, "password", "p", "", "MQTT Password")

	listenCmd.Flags().Bool("insecure", false, "Skip TLS certificate verification")
	viper.BindPFlag("insecure", listenCmd.Flags().Lookup("insecure"))

	listenCmd.Flags().String("ca-cert", "", "CA certificate file")
	viper.BindPFlag("ca-cert", listenCmd.Flags().Lookup("ca-cert"))
}

func tlsConfig() (*tls.Config, error) {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	if caCert, err := ioutil.ReadFile(viper.GetString("ca-cert")); err != nil {
		return nil, errors.New("Could not read CA certificate.")
	} else {
		certPool.AppendCertsFromPEM(caCert)
	}

	config := tls.Config{
		InsecureSkipVerify: viper.GetBool("insecure"),
		RootCAs:            certPool}

	return &config, nil
}
