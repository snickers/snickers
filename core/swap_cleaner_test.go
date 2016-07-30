package core_test

import (
	"io"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/core"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

var _ = Describe("Swap Cleaner", func() {
	var (
		dbInstance db.DatabaseInterface
		job        types.Job
		err        error
	)

	BeforeEach(func() {
		job = types.Job{
			ID:               "123",
			Source:           "http://source.here.mp4",
			Destination:      "s3://user@pass:/bucket/",
			Preset:           types.Preset{Name: "presetHere", Container: "mp4"},
			Status:           types.JobCreated,
			Details:          "",
			LocalSource:      "/tmp/123/src/KailuaBeach.mp4",
			LocalDestination: "/tmp/123/dst/KailuaBeach.webm",
		}

		dbInstance, err = db.GetDatabase()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		dbInstance.ClearDatabase()
	})

	Context("when calling", func() {
		JustBeforeEach(func() {
			dbInstance.StoreJob(job)
			Expect(os.MkdirAll("/tmp/123/src/", 0777)).To(Succeed())
			Expect(os.MkdirAll("/tmp/123/dst/", 0777)).To(Succeed())

			cp(job.LocalSource, "./videos/nyt.mp4")
			cp(job.LocalDestination, "./videos/nyt.mp4")

			Expect(job.LocalSource).To(BeAnExistingFile())
			Expect(job.LocalDestination).To(BeAnExistingFile())
		})

		It("should remove local source and local destination", func() {
			core.CleanSwap(job.ID)

			Expect(job.LocalSource).To(Not(BeAnExistingFile()))
			Expect(job.LocalDestination).To(Not(BeAnExistingFile()))
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
