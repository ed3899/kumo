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

// Returns a cobra command. The build command is used to build an AMI with ready to use tools.
func Build() *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Build an AMI with ready to use tools",
		Long:  `Build an AMI with ready to use tools.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			oopsBuilder := oops.
				Code("Build").
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
				Code("Build").
				In("cmd").
				Tags("cobra.Command").
				With("command", cmd.Name()).
				With("args", args)

			defer func() {
				if r := recover(); r != nil {
					err := oopsBuilder.Errorf("%v", r)
					log.Fatalf("panic: %+v", err)
				}
			}()

			_manager, err := manager.NewManager(iota.CloudIota(viper.GetString("cloud")), iota.Packer)
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to create new manager")

				panic(err)
			}

			err = _manager.SetCloudCredentials()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to set manager cloud credentials")

				panic(err)
			}
			defer _manager.UnsetCloudCredentials()

			err = _manager.SetPluginsPath()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to set plugins environment vars")

				panic(err)
			}
			defer _manager.UnsetPluginsEnvironmentVars()

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
			defer vars.Close()

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

			packer, err := binaries.NewPacker(_manager)
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to create new packer")

				panic(err)
			}

			err = _manager.GoToDirRun()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to chdir to manager dir")

				panic(err)
			}
			defer _manager.GoToDirInitial()

			err = packer.Init()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to init")

				panic(err)
			}

			err = packer.Build()
			if err != nil {
				err := oopsBuilder.
					Wrapf(err, "failed to build")

				panic(err)
			}
		},
	}
}
