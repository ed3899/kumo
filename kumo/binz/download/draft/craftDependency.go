package draft

import (
	"fmt"
	"log"

	"github.com/ed3899/kumo/host"
	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
)

type Dependency struct {
	Name           string
	Present        bool
	URL            string
	ZipPath        string
	ExtractionPath string
	ContentLength  int64
	DownloadBar    *mpb.Bar
	ZipBar         *mpb.Bar
}

// Return a dependency indicating whether the dependency is present. This only considers the presence
// of the downloaded zip file. If the downloaded zip file is present, the dependency is considered present.
func CraftDependency(name string) (dp *Dependency, err error) {
	depsDir := utils.GetDependenciesDirName()
	specs := host.GetSpecs()
	zipName := fmt.Sprintf("%s_%s_%s.zip", name, specs.OS, specs.ARCH)

	destinationZipPath, err := utils.CraftAbsolutePath(depsDir, zipName)
	if err != nil {
		msg := fmt.Sprintf("failed to get zip path for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	destinationExtractionPath, err := utils.CraftAbsolutePath(depsDir, name)
	if err != nil {
		msg := fmt.Sprintf("failed to get extraction path for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	url, err := utils.GetDependencyURL(name, specs)
	if err != nil {
		msg := fmt.Sprintf("failed to get url for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	contentLength, err := utils.GetContentLength(url)
	if err != nil {
		msg := fmt.Sprintf("failed to get content length for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	executablePath, err := utils.CraftAbsolutePath(depsDir, "packer", "packer.exe")
	if err != nil {
		msg := fmt.Sprintf("failed to get executable path for dependency: %v", name)
		err = errors.Wrap(err, msg)
		return nil, err
	}

	if utils.FileNotPresent(executablePath) {
		log.Printf("%s not present", name)
		return &Dependency{
			Name:           name,
			Present:        false,
			URL:            url,
			ExtractionPath: destinationExtractionPath,
			ZipPath:        destinationZipPath,
			ContentLength:  contentLength,
		}, nil
	}

	return &Dependency{
		Name:           name,
		Present:        true,
		URL:            url,
		ExtractionPath: destinationExtractionPath,
		ZipPath:        destinationZipPath,
		ContentLength:  contentLength,
	}, nil
}