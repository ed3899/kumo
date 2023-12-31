package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/common/constants"
	"github.com/ed3899/kumo/common/iota"
	"github.com/ed3899/kumo/manager/environment"
	"github.com/samber/oops"
)

// Creates a new manager. Used for cmd workflows
func NewManager(
	cloud iota.Cloud,
	tool iota.Tool,
) (*Manager, error) {
	oopsBuilder := oops.
		Code("NewManager").
		In("manager").
		Tags("Manager").
		With("cloud", cloud).
		With("tool", tool)

	currentExecutablePath, err := os.Executable()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get current executable path")

		return nil, err
	}
	currentExecutableDir := filepath.Dir(currentExecutablePath)

	currentWorkingDir, err := os.Getwd()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get current working directory")

		return nil, err
	}

	templatePath := func(templateName string) string {
		return filepath.Join(
			currentExecutableDir,
			iota.Templates.Name(),
			tool.Name(),
			templateName,
		)
	}

	pathToPackerManifest := filepath.Join(
		currentExecutableDir,
		iota.Packer.Name(),
		cloud.Name(),
		constants.PACKER_MANIFEST,
	)

	_environment, err := environment.NewEnvironment(
		tool,
		cloud,
		pathToPackerManifest,
	)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to create environment")

		return nil, err
	}

	terraformPath := func(fileName string) string {
		return filepath.Join(
			currentExecutableDir,
			iota.Terraform.Name(),
			cloud.Name(),
			fileName,
		)
	}

	return &Manager{
		Cloud: cloud.Iota(),
		Tool:  tool.Iota(),
		Path: &Path{
			Executable: filepath.Join(
				currentExecutableDir,
				iota.Dependencies.Name(),
				tool.Name(),
				fmt.Sprintf("%s.exe", tool.Name()),
			),
			Template: &Template{
				Merged: templatePath(constants.MERGED_TEMPLATE_NAME),
				Cloud:  templatePath(cloud.TemplateFiles().Cloud),
				Base:   templatePath(cloud.TemplateFiles().Base),
			},
			Vars: filepath.Join(
				currentExecutableDir,
				tool.Name(),
				cloud.Name(),
				tool.VarsName(),
			),
			Terraform: &Terraform{
				Lock:         terraformPath(constants.TERRAFORM_LOCK),
				State:        terraformPath(constants.TERRAFORM_STATE),
				Backup:       terraformPath(constants.TERRAFORM_BACKUP),
				IpFile:       terraformPath(constants.IP_FILE_NAME),
				IdentityFile: terraformPath(constants.KEY_NAME),
				SshConfig: filepath.Join(
					currentWorkingDir,
					constants.CONFIG_NAME,
				),
			},
			Dir: &Dir{
				Plugins: filepath.Join(
					currentExecutableDir,
					tool.Name(),
					cloud.Name(),
					tool.PluginDir(),
				),
				Initial: currentExecutableDir,
				Run: filepath.Join(
					currentExecutableDir,
					tool.Name(),
					cloud.Name(),
				),
			},
		},
		Environment: _environment,
	}, nil
}

type Manager struct {
	Cloud       iota.Cloud
	Tool        iota.Tool
	Path        *Path
	Environment any
}

type Path struct {
	Executable string
	Vars       string
	Terraform  *Terraform
	Template   *Template
	Dir        *Dir
}

type Terraform struct {
	Lock         string
	State        string
	Backup       string
	SshConfig    string
	IpFile       string
	IdentityFile string
}

type Template struct {
	Merged string
	Cloud  string
	Base   string
}

type Dir struct {
	Plugins string
	Initial string
	Run     string
}
