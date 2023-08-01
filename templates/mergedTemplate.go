package templates

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/ed3899/kumo/common/dirs"
	"github.com/ed3899/kumo/common/hashicorp_vars"
	"github.com/ed3899/kumo/common/templates"
	"github.com/ed3899/kumo/common/utils"
	"github.com/samber/oops"
	"go.uber.org/zap"
)

type MergedEnvironment[E templates.EnvironmentI] struct {
	General E
	Cloud   E
}

type MergedTemplate struct {
	instance    *template.Template
	absPath     string
	environment *MergedEnvironment[templates.EnvironmentI]
}

func NewMergedTemplate(generalTemplate, cloudTemplate templates.TemplateSingle) (packerMergedTemplate *MergedTemplate, err error) {
	var (
		oopsBuilder = oops.
				Code("new_packer_merged_template_failed").
				With("generalTemplate", generalTemplate.GetAbsPath()).
				With("cloudTemplate", cloudTemplate.GetAbsPath())

		mergedTemplateInstance     *template.Template
		absPathToTemplatesDir      string
		absPathToMergedTemplateDir string
		absPathToMergedTemplate    string
	)

	if generalTemplate.GetParentDirName() != cloudTemplate.GetParentDirName() {
		err = oopsBuilder.
			With("generalTemplate.GetParentDirName()", generalTemplate.GetParentDirName()).
			With("cloudTemplate.GetParentDirName()", cloudTemplate.GetParentDirName()).
			Errorf("generalTemplate and cloudTemplate must be in the same directory")
		return
	}

	if generalTemplate.GetEnvironment().IsNotValidEnvironment() || cloudTemplate.GetEnvironment().IsNotValidEnvironment() {
		err = oopsBuilder.
			Errorf("generalTemplate and cloudTemplate must have valid environments")
		return
	}

	if absPathToTemplatesDir, err = filepath.Abs(dirs.TEMPLATES_DIR_NAME); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while crafting absolute path to %s", dirs.TEMPLATES_DIR_NAME)
		return
	}

	absPathToMergedTemplateDir = filepath.Join(absPathToTemplatesDir, generalTemplate.GetParentDirName())

	if absPathToMergedTemplate, err = utils.MergeFilesTo(
		absPathToMergedTemplateDir,
		generalTemplate.GetAbsPath(),
		cloudTemplate.GetAbsPath(),
	); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while merging %s and %s to %s", generalTemplate.GetAbsPath(), cloudTemplate.GetAbsPath(), absPathToMergedTemplateDir)
		return
	}

	if mergedTemplateInstance, err = template.ParseFiles(absPathToMergedTemplate); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while parsing template %s", absPathToMergedTemplate)
		return
	}

	packerMergedTemplate = &MergedTemplate{
		instance: mergedTemplateInstance,
		absPath:  absPathToMergedTemplate,
		environment: &MergedEnvironment[templates.EnvironmentI]{
			General: generalTemplate.GetEnvironment(),
			Cloud:   cloudTemplate.GetEnvironment(),
		},
	}

	return
}

func (mt *MergedTemplate) GetAbsPath() (path string) {
	return mt.absPath
}

func (mt *MergedTemplate) GetName() (name string) {
	return mt.instance.Name()
}

func (mt *MergedTemplate) GetInstance() (instance *template.Template) {
	return mt.instance
}

func (mt *MergedTemplate) GetEnvironment() (environment *MergedEnvironment[templates.EnvironmentI]) {
	return mt.environment
}

func (mt *MergedTemplate) Remove() (err error) {
	var (
		logger, _ = zap.NewProduction()
	)

	defer logger.Sync()

	if err = os.Remove(mt.absPath); err != nil {
		logger.Sugar().Warnf("Failed to remove %s", mt.absPath)
		return
	}

	return
}

func (mt *MergedTemplate) ExecuteOn(hashicorpVars hashicorp_vars.HashicorpVarsI) (err error) {
	var (
		oopsBuilder = oops.
			Code("merged_template_execute_failed").
			With("hashicorpVars.GetFile().Name()", hashicorpVars.GetFile().Name())
	)

	if err = mt.instance.Execute(hashicorpVars.GetFile(), mt.environment); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while executing template: %s", mt.instance.Name())
		return
	}

	return
}