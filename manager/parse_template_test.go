package manager

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseTemplate", func() {
	Context("with a valid template", func() {
		var (
			content = `GIT_USERNAME = "{{.Base.Required.GIT_USERNAME}}"
			GIT_EMAIL = "{{.Base.Required.GIT_EMAIL}}"
			ANSIBLE_TAGS = [{{range $index, $element := .Base.Required.ANSIBLE_TAGS}}{{if $index}},{{end}}"{{$element}}"{{end}}]
			GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC = "{{.Base.Optional.GIT_HUB_PERSONAL_ACCESS_TOKEN_CLASSIC}}"`

			filePath string
			manager  *Manager
		)

		BeforeEach(func() {
			cwd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())

			filePath = filepath.Join(cwd, "merged")

			err = os.WriteFile(filePath, []byte(content), 0644)
			Expect(err).ToNot(HaveOccurred())

			manager = &Manager{
				Path: &Path{
					Template: &Template{
						Merged: filePath,
					},
				},
			}
		})

		AfterEach(func() {
			err := os.Remove(filePath)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should parse the template successfully", Label("unit"), func() {
			template, err := manager.ParseTemplate()
			Expect(err).ToNot(HaveOccurred())
			Expect(template).ToNot(BeNil())
		})
	})

	Context("with an invalid template", func() {
		var (
			manager *Manager
		)

		BeforeEach(func() {
			manager = &Manager{
				Path: &Path{
					Template: &Template{
						Merged: "invalid/path/to/template/file",
					},
				},
			}
		})

		It("should return an error", Label("unit"), func() {
			_, err := manager.ParseTemplate()
			Expect(err).To(HaveOccurred())
		})
	})
})