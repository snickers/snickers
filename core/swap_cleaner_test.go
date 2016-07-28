package core_test

import (
	"io"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/core"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers"
)

var _ = Describe("Swap Cleaner", func() {
	Context("when calling", func() {
		It("should remove local source and local destination", func() {
			dbInstance, _ := db.GetDatabase()
			dbInstance.ClearDatabase()

			exampleJob := snickers.Job{
				ID:               "123",
				Source:           "http://source.here.mp4",
				Destination:      "s3://user@pass:/bucket/",
				Preset:           snickers.Preset{Name: "presetHere", Container: "mp4"},
				Status:           snickers.JobCreated,
				Details:          "",
				LocalSource:      "/tmp/123/src/KailuaBeach.mp4",
				LocalDestination: "/tmp/123/dst/KailuaBeach.webm",
			}

			dbInstance.StoreJob(exampleJob)

			Expect(os.MkdirAll("/tmp/123/src/", 0777)).To(Succeed())
			Expect(os.MkdirAll("/tmp/123/dst/", 0777)).To(Succeed())

			cp(exampleJob.LocalSource, "./videos/nyt.mp4")
			cp(exampleJob.LocalDestination, "./videos/nyt.mp4")

			Expect(exampleJob.LocalSource).To(BeAnExistingFile())
			Expect(exampleJob.LocalDestination).To(BeAnExistingFile())

			core.CleanSwap(exampleJob.ID)

			Expect(exampleJob.LocalSource).To(Not(BeAnExistingFile()))
			Expect(exampleJob.LocalDestination).To(Not(BeAnExistingFile()))
		})
	})
})

func cp(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}
