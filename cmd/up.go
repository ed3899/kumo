package cmd

import (
	"log"
	"os"

	"github.com/ed3899/kumo/binaries"
	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/download"
	"github.com/ed3899/kumo/manager"
	"github.com/samber/oops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Returns a cobra command. The up command is used to deploy a cloud environment.
func Up() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Deploy your cloud environment",
		Long: `Deploy you cloud development environment. If no AMI is specified in the config file, Kumo will
		deploy the latest AMI built. It generates an SSH config file for you to easily SSH into your
		instances.`,
		Args: cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			oopsBuilder := oops.
				Code("Up").
				In("cmd").
				Tags("Cobra", "PreRun")

			cwd, err := os.Getwd()
			if err != nil {
				log.Fatalf(
					"%+v",
					oopsBuilder.
						Wrapf(err, "Error occurred while getting current working directory"),
				)
			}

			viper.SetConfigName("kumo.config")
			viper.SetConfigType("yaml")
			viper.AddConfigPath(cwd)

			err = viper.ReadInConfig()
			if err != nil {
				log.Fatalf(
					"%+v",
					oopsBuilder.
						Wrapf(err, "Error occurred while reading config file. Make sure a kumo.config.yaml file exists in the current working directory"),
				)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			oopsBuilder := oops.
				Code("Up").
				In("cmd").
				Tags("Cobra").
				Tags("Run").
				With("command", *cmd).
				With("args", args)

			defer func() {
				if r := recover(); r != nil {
					err := oopsBuilder.Errorf("%v", r)
					log.Fatalf("panic: %+v", err)
				}
			}()

			_manager, err := manager.NewManager(iota.CloudIota(viper.GetString("cloud")), iota.Terraform)
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to create new manager")

				panic(err)
			}

			err = _manager.CreateTemplate()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to create template")

				panic(err)
			}

			template, err := _manager.ParseTemplate()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to parse template")

				panic(err)
			}
			defer _manager.DeleteTemplate()

			vars, err := _manager.CreateVars()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to create vars")

				panic(err)
			}

			err = template.Execute(vars, _manager.Environment)
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to execute template")

				panic(err)
			}

			if !_manager.ToolExecutableExists() {
				_download, err := download.NewDownload(_manager)
				if err != nil {
					err := oopsBuilder.
						Wrapf(err, "failed to create new download")

					panic(err)
				}

				defer _download.RemoveZip()

				err = _download.DownloadAndShowProgress()
				if err != nil {
					err := oopsBuilder.
						Wrapf(err, "failed to download")

					panic(err)
				}

				err = _download.ExtractAndShowProgress()
				if err != nil {
					err := oopsBuilder.
						Wrapf(err, "failed to extract")

					panic(err)
				}

				_download.ProgressShutdown()
			}

			terraform, err := binaries.NewTerraform(_manager)
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to create new terraform")

				panic(err)
			}

			err = _manager.GoToDirRun()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to chdir to manager dir")

				panic(err)
			}
			defer _manager.GoToDirInitial()

			err = terraform.Init()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to init")

				panic(err)
			}

			err = terraform.Apply()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to apply")

				panic(err)
			}

			err = _manager.CreateSshConfig()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to generate ssh config")

				panic(err)
			}
		},
	}
}
