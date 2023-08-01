package general

import (
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/templates"
	"github.com/ed3899/kumo/common/tool"
	"github.com/ed3899/kumo/common/utils"
	"github.com/samber/oops"
	"go.uber.org/zap"
)

type Template struct {
	instance      *template.Template
	absPath       string
	parentDirName string
	environment   templates.EnvironmentI
}

func NewTemplate() (newTemplate *Template, err error) {
	var (
		oopsBuilder = oops.
				Code("new_template_failed")
		logger, _                    = zap.NewProduction()
		templatesDirName             = dirs.TEMPLATES_DIR_NAME
		terraformDirName             = tool.TERRAFORM_NAME
		generalDirName               = dirs.GENERAL_DIR_NAME
		terraformGeneralTemplateName = templates.TERRAFORM_GENERAL_TEMPLATE_NAME
		terraformDefaultAllowedIp    = tool.TERRAFORM_DEFAULT_ALLOWED_IP

		absPathToTemplate string
		instance          *template.Template
		publicIp          string
		pickedIp          string
	)

	defer logger.Sync()

	if absPathToTemplate, err = filepath.Abs(filepath.Join(templatesDirName, terraformDirName, generalDirName, terraformGeneralTemplateName)); err != nil {
		err = oopsBuilder.
			With("templatesDirName", templatesDirName).
			With("terraformDirName", terraformDirName).
			With("generalDirName", generalDirName).
			Wrapf(err, "Error occurred while crafting absolute path to %s", terraformGeneralTemplateName)
		return
	}

	if instance, err = template.ParseFiles(absPathToTemplate); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPathToTemplate)
		return
	}

	if publicIp, err = utils.GetPublicIp(); err != nil {
		logger.Sugar().Warnf("Failed to get public IP address. Defaulting to %s", terraformDefaultAllowedIp)
		pickedIp = terraformDefaultAllowedIp
	} else {
		pickedIp = publicIp
	}

	newTemplate = &Template{
		instance:      instance,
		absPath:       absPathToTemplate,
		parentDirName: terraformDirName,
		environment: &Environment{
			ALLOWED_IP: utils.MaskIp(pickedIp, 32),
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