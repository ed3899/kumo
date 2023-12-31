package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/samber/oops"
)

// Unzips the zip file at the given path to the given destination path.
// Sends the number of bytes unzipped to the given channel.
//
// Example:
//	("/home/dev/packer_1.7.4_windows_amd64.zip", "/home/dev/unzip", bytesUnzippedChan) -> nil
func Unzip(
	pathToZip, pathToUnzip string,
	bytesUnzipped chan<- int,
) error {
	oopsBuilder := oops.
		Code("Unzip").
		In("utils").
		In("zip").
		With("pathToZip", pathToZip).
		With("pathToUnzip", pathToUnzip).
		With("bytesUnzipped", bytesUnzipped)

	// Open the zip file and defer closing it
	reader, err := zip.OpenReader(pathToZip)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to open zip file: %s", pathToZip)
		return err
	}
	defer reader.Close()

	unzipGroup := &sync.WaitGroup{}
	errChan := make(chan error, len(reader.File))

	// Unzip each file concurrently
	for _, zipFile := range reader.File {
		unzipGroup.Add(1)
		go func(zf *zip.File) {
			defer unzipGroup.Done()

			var (
				bytesCopied int64
			)

			if bytesCopied, err = UnzipFile(zf, pathToUnzip); err != nil {
				err = oopsBuilder.
					With("bytesCopied", bytesCopied).
					With("zipFile", zf.Name).
					With("extractToPath", pathToUnzip).
					Wrapf(err, "failed to unzip file: %s", zf.Name)
				errChan <- err
				return
			}

			bytesUnzipped <- int(bytesCopied)
		}(zipFile)
	}

	// Wait for all files to be unzipped
	go func() {
		unzipGroup.Wait()
		close(errChan)
		close(bytesUnzipped)
	}()

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// Unzips a file from a zip file to a destination path.
// The destination path is created if it doesn't exist.
// Returns the number of bytes copied.
//
// Example:
//	(zipFile, "/home/dev/unzip") -> 1234, nil
func UnzipFile(
	zf *zip.File,
	extractToPath string,
) (int64, error) {
	oopsBuilder := oops.
		Code("UnzipFile").
		In("utils").
		In("zip").
		With("zf", zf).
		With("extractToPath", extractToPath)

	// Check if file path is not vulnerable to Zip Slip
	filePath := filepath.Join(extractToPath, zf.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(extractToPath)+string(os.PathSeparator)) {
		err := oopsBuilder.
			Errorf("illegal file path: %s", filePath)
		return -1, err
	}

	// Check if file is a directory
	if zf.FileInfo().IsDir() {
		err := oopsBuilder.
			Errorf("is a directory: %s", zf.Name)
		return -1, err
	}

	// Create directory tree
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to create directory tree for: %s", filePath)
		return -1, err
	}

	// Create a destination file for unzipped content and defer closing it
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zf.Mode())
	if err != nil {
		err := oopsBuilder.
			With("filePath", filePath).
			Wrapf(err, "failed to create destination file: %s", filePath)
		return -1, err
	}
	defer destinationFile.Close()

	// Unzip the content of a file and copy it to the destination file. Defer closing the zipped file
	zippedFile, err := zf.Open()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to open zipped file: %s", zf.Name)
		return -1, err
	}
	defer zippedFile.Close()

	bytesCopied, err := io.Copy(destinationFile, zippedFile)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to copy zipped file %#v to destination file: %s", zippedFile, destinationFile.Name())
		return -1, err
	}

	return bytesCopied, nil
}
