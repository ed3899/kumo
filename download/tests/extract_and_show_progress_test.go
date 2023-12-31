package tests

import (
	"archive/zip"
	"os"
	"path/filepath"
	"sync"

	"github.com/ed3899/kumo/download"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vbauerster/mpb/v8"
)

var _ = Describe("ExtractAndShowProgress", func() {
	Context("with a valid zip", func() {
		var (
			_download *download.Download
			ziPath   string
			zipFile  *os.File
			exePath  string
		)

		BeforeEach(func() {
			cwd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())

			ziPath = filepath.Join(cwd, "mock.zip")

			zipFile, err = os.Create(ziPath)
			Expect(err).ToNot(HaveOccurred())
			defer zipFile.Close()

			// Create a new zip writer
			zipWriter := zip.NewWriter(zipFile)
			defer zipWriter.Close()

			// Define the file content
			fileContent := []byte("This is mock content.")

			// Create a file in the zip archive
			exeName := "mock.txt"
			file, err := zipWriter.Create(exeName)
			Expect(err).ToNot(HaveOccurred())

			// Write content to the file in the zip
			_, err = file.Write(fileContent)
			Expect(err).ToNot(HaveOccurred())

			exePath = filepath.Join(cwd, exeName)

			_download = &download.Download{
				Path: &download.Path{
					Zip:        zipFile.Name(),
					Executable: exePath,
				},
				Progress: mpb.New(mpb.WithWaitGroup(&sync.WaitGroup{}), mpb.WithAutoRefresh(), mpb.WithWidth(0)),
				Bar:      &download.Bar{},
			}
		})

		AfterEach(func() {
			err := os.Remove(zipFile.Name())
			Expect(err).ToNot(HaveOccurred())

			err = os.Remove(exePath)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should successfully extract and show progress", func() {
			err := _download.ExtractAndShowProgress()
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("with an invalid zip", func() {
		var (
			_download *download.Download
		)

		BeforeEach(func() {
			_download = &download.Download{
				Path: &download.Path{
					Zip: "invalid_zip",
				},
				Progress: mpb.New(mpb.WithWaitGroup(&sync.WaitGroup{}), mpb.WithAutoRefresh(), mpb.WithWidth(0)),
				Bar:      &download.Bar{},
			}
		})

		It("should handle extract error", func() {
			err := _download.ExtractAndShowProgress()
			Expect(err).To(HaveOccurred())
		})
	})
})
