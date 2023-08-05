package tool

import (
	constants "github.com/ed3899/kumo/0_constants"
	"github.com/samber/oops"
)

type Tool struct {
	Name              string
	Version           string
}

func NewTool(kind constants.ToolKind) (tool Tool, err error) {
	var (
		oopsBuilder = oops.
			Code("new_tool_setup_failed").
			With("tool", kind)
	)

	switch kind {
	case constants.Packer:
		tool = Tool{
			Name:    constants.PACKER,
			Version: constants.PACKER_VERSION,
		}

	case constants.Terraform:
		tool = Tool{
			Name:    constants.TERRAFORM,
			Version: constants.TERRAFORM_VERSION,
		}

	default:
		err = oopsBuilder.
			Errorf("Unknown tool kind: %d", kind)
		return

	}

	return
}