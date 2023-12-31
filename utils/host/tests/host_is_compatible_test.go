package tests

import (
	"github.com/ed3899/kumo/utils/host"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HostIsCompatible", func() {
	Context("on windows", func() {
		var windows = "windows"

		It("should return true for compatible architectures", Label("unit"), func() {
			Expect(host.HostIsCompatible(windows, "386")).To(BeTrue())
			Expect(host.HostIsCompatible(windows, "amd64")).To(BeTrue())
		})

		It("should return false for incompatible architectures", Label("unit"), func() {
			Expect(host.HostIsCompatible(windows, "arm64")).To(BeFalse())
		})
	})

	Context("on darwin", func() {
		var darwin = "darwin"

		It("should return true for compatible architectures", Label("unit"), func() {
			Expect(host.HostIsCompatible(darwin, "amd64")).To(BeTrue())
			Expect(host.HostIsCompatible(darwin, "arm64")).To(BeTrue())
		})

		It("should return false for incompatible architectures", Label("unit"), func() {
			Expect(host.HostIsCompatible(darwin, "386")).To(BeFalse())
		})
	})

	Context("on incompatible platforms", func() {
		It("should return false for incompatible platforms", Label("unit"), func() {
			Expect(host.HostIsCompatible("linux", "amd64")).To(BeFalse())
		})
	})
})
