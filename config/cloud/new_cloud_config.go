package cloud

import (
	"path/filepath"

	"github.com/ed3899/kumo/common/interfaces"
	"github.com/ed3899/kumo/constants"
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

func NewCloudConfig(
	options ...Option,
) (
	cloud *CloudConfig,
	err error,
) {
	var (
		oopsBuilder = oops.
				Code("NewCloud").
				With("opts", options)

		opt Option
	)

	cloud = &CloudConfig{}
	for _, opt = range options {
		if err = opt(cloud); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Option %v", opt)
			return
		}
	}

	return
}

func WithKind(
	cloudFromConfig string,
) (
	option Option,
) {
	var (
		oopsBuilder = oops.
			Code("WithKind").
			With("cloudFromConfig", cloudFromConfig)
	)

	option = func(cloud *CloudConfig) (err error) {
		switch cloudFromConfig {
		case constants.AWS:
			cloud.Kind = constants.Aws
		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud %s is not supported", cloudFromConfig)
			return
		}

		return
	}

	return
}

func WithName(
	cloudFromConfig string,
) (
	option Option,
) {
	var (
		oopsBuilder = oops.
			Code("WithName").
			With("cloudFromConfig", cloudFromConfig)
	)

	option = func(cloud *CloudConfig) (err error) {
		switch cloudFromConfig {
		case constants.AWS:
			cloud.Name = constants.AWS
		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud %s is not supported", cloudFromConfig)
			return
		}

		return
	}

	return
}

func WithCredentials(
	cloudFromConfig string,
) (
	option Option,
) {
	var (
		oopsBuilder = oops.
			Code("WithCredentials").
			With("cloudFromConfig", cloudFromConfig)
	)

	option = func(cloud *CloudConfig) (err error) {
		switch cloudFromConfig {
		case constants.AWS:
			cloud.Credentials = AwsCredentials{
				AccessKeyId:     viper.GetString("AWS.AccessKeyId"),
				SecretAccessKey: viper.GetString("AWS.SecretAccessKey"),
			}
		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud %s is not supported", cloudFromConfig)
			return
		}

		return
	}

	return
}

func WithPackerManifestPath(
	cloudFromConfig,
	kumoExecAbsPath string,
) (
	option Option,
) {
	var (
		oopsBuilder = oops.
			Code("WithPackerManifestPath").
			With("cloudFromConfig", cloudFromConfig)
	)

	option = func(cloud *CloudConfig) (err error) {
		switch cloudFromConfig {
		case constants.AWS:
			cloud.PackerManifestPath = filepath.Join(kumoExecAbsPath, constants.PACKER, constants.AWS, constants.PACKER_MANIFEST)
		default:
			err = oopsBuilder.
				Wrapf(err, "Cloud %s is not supported", cloudFromConfig)
			return
		}

		return
	}

	return
}

type CloudConfig struct {
	Kind               constants.CloudKind
	Name               string
	Credentials        interfaces.Credentials
	PackerManifestPath string
}

type Option func(*CloudConfig) error
