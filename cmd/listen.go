// Copyright © 2017 Philippe Hässig <phil@neckhair.ch>

package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/neckhair/owntracks-eventr/listener"
)

var config = listener.Configuration{}

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for events and write them into a file",
	Long: `Listen for events and write them into a file line by line.

A password for MQTT can be provided in an environment variable named MQTT_PASSWORD.

    MQTT_PASSWORD=secret owntracks-eventr listen -u eventr
`,

	PreRun: func(cmd *cobra.Command, args []string) {
		config.Url = viper.GetString("url")
		config.Filename = viper.GetString("output")
		config.Username = viper.GetString("username")
		config.Password = viper.GetString("password")

		log.Printf("--> Listening for MQTT events\n")
		log.Printf("Server:  %s\n", config.Url)
		log.Printf("Output:  %s\n", config.Filename)
		log.Printf("Logfile: %s\n\n", viper.GetString("LogFile"))
	},

	Run: func(cmd *cobra.Command, args []string) {
		config.Password = viper.GetString("password")
		listener := listener.NewListener(&config)

		var err error
		if listener.TLSConfig, err = tlsConfig(); err != nil {
			log.Fatalln(err)
		}

		if err := listener.Start(); err != nil {
			log.Fatalf("Could not connect to MQTT server. %s\n", err)
		}
		defer listener.Stop()

		waitForQuit()
	},
}

func init() {
	RootCmd.AddCommand(listenCmd)

	listenCmd.Flags().StringP("url", "", "tls://localhost:8883", "Connection string")
	listenCmd.Flags().StringP("output", "o", "eventlog.txt", "Path to destination file")
	listenCmd.Flags().StringP("username", "u", "", "MQTT Username")
	listenCmd.Flags().Bool("insecure", false, "Skip TLS certificate verification")
	listenCmd.Flags().String("ca-cert", "", "CA certificate file")

	viper.BindPFlag("url", listenCmd.Flags().Lookup("url"))
	viper.BindPFlag("output", listenCmd.Flags().Lookup("output"))
	viper.BindPFlag("username", listenCmd.Flags().Lookup("username"))
	viper.BindPFlag("insecure", listenCmd.Flags().Lookup("insecure"))
	viper.BindPFlag("ca-cert", listenCmd.Flags().Lookup("ca-cert"))

	viper.BindEnv("password", "MQTT_PASSWORD")
}

func tlsConfig() (*tls.Config, error) {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	if caCertPath := viper.GetString("ca-cert"); caCertPath != "" {
		if caCert, err := ioutil.ReadFile(viper.GetString("ca-cert")); err != nil {
			return nil, errors.New("Could not read CA certificate.")
		} else {
			certPool.AppendCertsFromPEM(caCert)
		}
	}

	config := tls.Config{
		InsecureSkipVerify: viper.GetBool("insecure"),
		RootCAs:            certPool}

	return &config, nil
}

func waitForQuit() {
	var endWaiter sync.WaitGroup
	endWaiter.Add(1)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChannel
		endWaiter.Done()
	}()

	endWaiter.Wait()
}
