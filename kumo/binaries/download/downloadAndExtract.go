package download

import (
	"os"
	"path/filepath"

	"github.com/ed3899/kumo/binaries"
	"github.com/samber/oops"
	"github.com/vbauerster/mpb/v8"
)

func DownloadAndExtract(z binaries.ZipI, extractToAbsPathDir string) (err error) {
	var (
		absPathToZipDir = filepath.Dir(z.GetPath())
		progress        = mpb.New(mpb.WithWidth(64), mpb.WithAutoRefresh())
		oopsBuilder     = oops.
				Code("download_and_extract_failed").
				With("extractToAbsPathDir", extractToAbsPathDir).
				With("z.GetName()", z.GetName()).
				With("z.GetPath()", z.GetPath())
	)

	// Start with a clean slate
	if err = os.RemoveAll(extractToAbsPathDir); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", extractToAbsPathDir)
		return
	}

	if err = os.RemoveAll(absPathToZipDir); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", absPathToZipDir)
		return
	}

	// Download
	if err = DownloadAndShowProgress(z, progress); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while downloading %s", z.GetName())
		return
	}

	// Extract
	if err = ExtractAndShowProgress(z, extractToAbsPathDir, progress); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while extracting %s", z.GetName())
		return
	}

	progress.Shutdown()

	// Remove zip
	if err = z.Remove(); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while removing %s", z.GetName())
		return
	}

	return
}