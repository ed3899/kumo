package utils

import (
	"fmt"
)

func GetLatestTerraformVersion() string {
	return "1.5.3"
}

func CreateHashicorpURL(name, version, os, arch string) (string) {
	return fmt.Sprintf("https://releases.hashicorp.com/%s/%s/%s_%s_%s_%s.zip", name, version, name, version, os, arch)
}