package utils

import (
	"fmt"
)

type CreateHashicorpURLF func(name, version, os, arch string) (url string)

func CreateHashicorpURL(name, version, os, arch string) string {
	return fmt.Sprintf(
		"https://releases.hashicorp.com/%s/%s/%s_%s_%s_%s.zip",
		name,
		version,
		name,
		version,
		os,
		arch,
	)
}
