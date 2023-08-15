package manager

import (
	"os"

	"github.com/ed3899/kumo/common/iota"
	"github.com/samber/oops"
)

func SetManagerCloudCredentials(
	manager *Manager,
) error {
	oopsBuilder := oops.
		In("manager").
		Tags("Manager").
		Code("SetCloudCredentials")

	switch manager.Cloud {
	case iota.Aws:
		for key, value := range awsCredentials {
			if err := os.Setenv(key, value); err != nil {
				return oopsBuilder.
					With("cloudName", manager.Cloud.Name).
					Wrapf(err, "failed to set environment variable %s to %s", key, value)
			}
		}

	default:
		return oopsBuilder.
			With("cloudName", manager.Cloud.Name).
			Errorf("unknown cloud: %#v", manager.Cloud)
	}

	return nil
}