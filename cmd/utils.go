package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

func addTargetFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&target.DomainName, "domainName", "d", "", "The domain name you wish to flood")
	cmd.Flags().IPVar(&target.IPAddress, "ip", nil, "The IP address you wish to target")
	cmd.Flags().Uint16VarP(&target.DstPort, "port", "p", 0, "The port you wish to target")
	cmd.Flags().StringVarP(&target.InterfaceName, "networkInterface", "i", "", "The network interface you wish to target")
}
