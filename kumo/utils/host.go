package utils

import "runtime"

func GetCurrentHostSpecs() (os, arch string) {
	return runtime.GOOS, runtime.GOARCH
}

func HostIsCompatible() (compatible bool) {
	switch runtime.GOOS {
	case "windows":
		switch runtime.GOARCH {
		case "386":
			compatible = true
		case "amd64":
			compatible = true
		default:
			compatible = false
		}
	default:
		compatible = false
	}
	return
}

func HostIsNotCompatible() bool {
	return !HostIsCompatible()
}