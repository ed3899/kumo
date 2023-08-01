package cmd

import (
	"log"

	"github.com/ed3899/kumo/workflows"
	"github.com/samber/oops"
	"github.com/spf13/cobra"
)

func DestroyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "destroy",
		Short: "Destroy your cloud environment",
		Long:  `Destroy your last deployed cloud environment. Doesn't destroy the AMI. It will also remove the SSH config file.`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				oopsBuilder = oops.Code("destroy_command_failed").
					With("command", cmd.Name()).
					With("args", args)
			)

			if err := workflows.Destroy(); err != nil {
				log.Fatalf(
					"%+v",
					oopsBuilder.
						Wrapf(err, "Error occurred running destroy workflow"),
				)
			}
		},
	}
}