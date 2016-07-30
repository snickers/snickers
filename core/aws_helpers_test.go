package core_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/snickers/snickers/core"
	"github.com/snickers/snickers/db"
)

var _ = Describe("AWS Helpers", func() {
	var (
		dbInstance  db.DatabaseInterface
		destination string
	)

	BeforeEach(func() {
		dbInstance, _ = db.GetDatabase()
		destination = "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT/HERE.mp4"
	})

	AfterEach(func() {
		dbInstance.ClearDatabase()
	})

	It("Should get bucket from URL Destination", func() {
		bucket, err := core.GetAWSBucket(destination)
		Expect(err).NotTo(HaveOccurred())
		Expect(bucket).To(Equal("BUCKET"))
	})

	It("Should set credentials from URL Destination", func() {
		Expect(core.SetAWSCredentials(destination)).To(Succeed())
		Expect(os.Getenv("AWS_ACCESS_KEY_ID")).To(Equal("AWSKEY"))
		Expect(os.Getenv("AWS_SECRET_ACCESS_KEY")).To(Equal("AWSSECRET"))
	})

	It("Should get path and filename from URL Destination", func() {
		key, err := core.GetAWSKey(destination)
		Expect(err).NotTo(HaveOccurred())
		Expect(key).To(Equal("/OBJECT/HERE.mp4"))
	})
})
