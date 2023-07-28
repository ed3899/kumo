package cloud

import (
	"github.com/samber/oops"
	"github.com/spf13/viper"
)

type Credentials interface {
	Set() error
	Unset() error
}

type CloudSetup struct {
	cloudName string

	Credentials Credentials
}

func NewCloudSetup(cloud string, tool Tool) (cloudSetup *CloudSetup, err error) {
	var (
		oopsBuilder = oops.
				Code("new_cloud_deployment_failed").
				With("cloud", cloud)

		cloudName   string
		credentials Credentials
	)

	switch cloud {
	case "aws":
		cloudName = "aws"
		credentials = &AwsCredentials{
			AccessKeyId:     viper.GetString("AWS.AccessKeyId"),
			SecretAccessKey: viper.GetString("AWS.SecretAccessKey"),
		}

	default:
		err = oopsBuilder.
			Errorf("Cloud '%v' not supported", cloud)
		return
	}

	cloudSetup = &CloudSetup{
		cloudName: cloudName,

		Credentials: credentials,
	}

	return
}

func (cs *CloudSetup) GetCloudName() (cloudName string) {
	return cs.cloudName
}
