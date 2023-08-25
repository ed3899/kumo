package tests

import (
	"github.com/ed3899/kumo/manager/environment"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewPackerEnvironment", func() {
	Context("when the cloud is unknown", func() {
		It("returns an error", func() {
			_, err := environment.NewPackerEnvironment(
				-1,
			)
			Expect(err).To(HaveOccurred())
		})
	})
})
