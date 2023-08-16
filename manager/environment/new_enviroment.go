package environment

import (
	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func NewEnvironment(tool iota.Tool, cloud iota.Cloud, pathToPackerManifest string) (*Environment[any], error) {
	oopsBuilder := oops.
		Code("NewEnvironment").
		With("tool", tool)

	switch tool {
	case iota.Packer:
		packerEnvironment, err := NewPackerEnvironment(cloud)
		if err != nil {
			return nil, oopsBuilder.
				With("cloud", cloud).
				Wrapf(err, "failed to create packer environment")
		}

		return &Environment[any]{
			General: packerEnvironment.General,
			Cloud:   packerEnvironment.Cloud,
		}, nil

	case iota.Terraform:
		terraformEnvironment, err := NewTerraformEnvironment(pathToPackerManifest, cloud)
		if err != nil {
			return nil, oopsBuilder.
				With("cloud", cloud).
				Wrapf(err, "failed to create terraform environment")
		}

		return &Environment[any]{
			General: terraformEnvironment.General,
			Cloud:   terraformEnvironment.Cloud,
		}, nil

	default:
		return nil, oopsBuilder.
			Errorf("unknown tool: %v", tool)
	}
}

type Environment[T any] struct {
	General T
	Cloud   any
}