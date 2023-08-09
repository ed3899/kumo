package cmd

import (
	"log"

	"github.com/ed3899/kumo/config/file"
	"github.com/samber/oops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Long: `🌩️ Kumo: Your quick and easy cloud development environment.`,
}

func init() {
	var (
		kumoConfigFile = file.NewKumoConfig(
			file.WithName("kumo.config"),
			file.WithType("yaml"),
			file.WithPath("."),
		).SetConfigName(viper.SetConfigName).
			SetConfigType(viper.SetConfigType).
			AddConfigPath(viper.AddConfigPath)

		err error
	)

	if err = kumoConfigFile.ReadInConfig(viper.ReadInConfig); err != nil {
		log.Fatalf(
			"%+v",
			oops.Code("cmd-root.go-init").
				Wrapf(err, "Error occurred while reading kumo config"),
		)
	}

	// Assemble commands
	rootCmd.AddCommand(*GetCommands()...)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf(
			"%+v",
			oops.Code("root_cmd_execute_failed").
				Wrapf(err, "Error occurred while running kumo"),
		)
	}
}
