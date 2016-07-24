package snickers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/core"
	"github.com/snickers/snickers/types"
)

var _ = Describe("Swap Cleaner", func() {
	Context("when calling", func() {
		It("should remove local source and local destination", func() {
			exampleJob := types.Job{
				ID:               "123",
				Source:           "http://source.here.mp4",
				Destination:      "s3://user@pass:/bucket/",
				Preset:           types.Preset{Name: "presetHere", Container: "mp4"},
				Status:           types.JobCreated,
				Details:          "",
				LocalSource:      "/tmp/123/src/KailuaBeach.mp4",
				LocalDestination: "/tmp/123/dst/KailuaBeach.webm",
			}

			err := core.CleanSwap(exampleJob.ID)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
