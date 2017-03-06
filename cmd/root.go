// Copyright © 2017 Levi Gross <levi@levigross.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/levigross/mabul/base"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var target base.Target
var loggingLevel, cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mabul",
	Short: "Mabul is a program designed as a test suite for DDoS mitigation programs",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mabul.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.PersistentFlags().StringVar(&loggingLevel, "logLevel", "info", "The level of logging you wish to have")
	RootCmd.PersistentFlags().StringVarP(&target.DomainName, "domainName", "d", "", "The domain name you wish to flood")
	RootCmd.PersistentFlags().IPVar(&target.IPAddress, "ip", nil, "The IP address you wish to target")
	RootCmd.PersistentFlags().IntVarP(&target.Port, "port", "p", 0, "The port you wish to target")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".mabul") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func parseLogLevel() zapcore.Level {
	switch strings.ToLower(loggingLevel) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	default:
		return zap.WarnLevel
	}
}

func setupLogging() *zap.SugaredLogger {
	loggingConfig := zap.NewDevelopmentConfig()
	loggingConfig.Level.SetLevel(parseLogLevel())
	logger, err := loggingConfig.Build()
	if err != nil {
		fmt.Println("Unable to create logger", err)
		os.Exit(-1)
	}
	return logger.Sugar()
}
