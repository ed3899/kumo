package ip

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetPublicIp", func() {
	It("should retrieve the public IP address", Label("integration"), func() {
		ip, err := GetPublicIp()
		Expect(err).NotTo(HaveOccurred())
		Expect(ip).To(MatchRegexp("\\b(?:\\d{1,3}\\.){3}\\d{1,3}\\b"))
	})
})
