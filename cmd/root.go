// Copyright © 2017 Philippe Hässig <phil@neckhair.ch>

package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile        string
	logFileHandler *os.File
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "owntracks-eventr",
	Short: "Listens for MQTT events from Owntrack and logs them into a file",
	Long: `owntracks-eventr listens on MQTT for events from Owntrack. It writes them
into a log file where you can calculate times spent at a location.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logFileHandler, err := os.OpenFile(viper.GetString("LogFile"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Cannot open logfile!")
			os.Exit(1)
		}

		log.SetOutput(logFileHandler)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		logFileHandler.Close()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("version") {
			fmt.Printf("owntracks-eventr %s (%s-%s)\n", Version, runtime.GOOS, runtime.GOARCH)
		} else {
			cmd.Usage()
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringP("logFile", "l", "/dev/stdout", "Log File path")

	RootCmd.Flags().BoolP("version", "v", false, "Show version number and quit")

	viper.BindPFlag("logFile", RootCmd.PersistentFlags().Lookup("logFile"))
	viper.BindPFlag("version", RootCmd.Flags().Lookup("version"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".owntracks-eventr") // name of config file (without extension)
	viper.AddConfigPath("$HOME")             // adding home directory as first search path

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
