package aws

import (
	"path/filepath"
	"text/template"

	common_cloud_constants "github.com/ed3899/kumo/common/cloud/constants"
	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/ssh"

	common_templates_interfaces "github.com/ed3899/kumo/common/templates/interfaces"
	common_templates_constants "github.com/ed3899/kumo/common/templates/constants"
	common_tool_constants "github.com/ed3899/kumo/common/tool/constants"
	"github.com/ed3899/kumo/common/utils"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Template struct {
	instance      *template.Template
	absPath       string
	parentDirName string
	environment   common_templates_interfaces.Environment
}

func NewTemplate(packerManifest common_templates_interfaces.PackerManifest) (newTemplate *Template, err error) {
	var (
		oopsBuilder = oops.
				Code("new_template_failed")
		templatesDirName         = dirs.TEMPLATES_DIR_NAME
		terraformDirName         = common_tool_constants.TERRAFORM_NAME
		awsDirName               = common_cloud_constants.AWS_NAME
		terraformAwsTemplateName = common_templates_constants.TERRAFORM_AWS_TEMPLATE_NAME

		absPathToTemplate string
		instance          *template.Template
		pickedAmiId       string
	)

	if absPathToTemplate, err = filepath.Abs(filepath.Join(templatesDirName, terraformDirName, awsDirName, terraformAwsTemplateName)); err != nil {
		err = oopsBuilder.
			With("templatesDirName", templatesDirName).
			With("terraformDirName", terraformDirName).
			With("awsDirName", awsDirName).
			Wrapf(err, "Error occurred while crafting absolute path to %s", terraformAwsTemplateName)
		return
	}

	if instance, err = template.ParseFiles(absPathToTemplate); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPathToTemplate)
		return
	}

	if pickedAmiId, err = utils.PickAmiIdToBeUsed(packerManifest.LastBuiltAmiId(), viper.GetString("AMI.IdToBeUsed")); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while picking ami id to be used")
		return
	}

	newTemplate = &Template{
		instance:      instance,
		absPath:       absPathToTemplate,
		parentDirName: terraformDirName,
		environment: &Environment{
			Required: &Required{
				AWS_REGION:        viper.GetString("AWS.Region"),
				AWS_INSTANCE_TYPE: viper.GetString("AWS.EC2.Instance.Type"),
				AMI_ID:            pickedAmiId,
				KEY_NAME:          ssh.KEY_NAME,
				SSH_PORT:          ssh.SSH_PORT,
				IP_FILE_NAME:      ssh.IP_FILE_NAME,
				USERNAME:          viper.GetString("AMI.User"),
			},
			Optional: &Optional{
				AWS_EC2_INSTANCE_VOLUME_TYPE: viper.GetString("AWS.EC2.Volume.Type"),
				AWS_EC2_INSTANCE_VOLUME_SIZE: viper.GetInt("AWS.EC2.Volume.Size"),
			},
		},
	}

	return
}

func (t *Template) GetAbsPath() (absPath string) {
	return t.absPath
}

func (t *Template) GetParentDirName() (dir string) {
	return t.parentDirName
}

func (t *Template) GetInstance() (instance *template.Template) {
	return t.instance
}

func (t *Template) GetEnvironment() (environment templates.EnvironmentI) {
	return t.environment
}