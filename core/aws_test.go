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
		dbInstance db.DatabaseInterface
	)

	BeforeEach(func() {
		dbInstance, _ = db.GetDatabase()
		dbInstance.ClearDatabase()
	})

	It("Should get bucket from URL Destination", func() {
		destination := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT"
		bucket, _ := core.GetAWSBucket(destination)
		Expect(bucket).To(Equal("BUCKET"))
	})

	It("Should set credentials from URL Destination", func() {
		destination := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT"
		core.SetAWSCredentials(destination)
		Expect(os.Getenv("AWS_ACCESS_KEY_ID")).To(Equal("AWSKEY"))
		Expect(os.Getenv("AWS_SECRET_ACCESS_KEY")).To(Equal("AWSSECRET"))
	})

	It("Should get path and filename from URL Destination", func() {
		destination := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT/HERE.mp4"
		key, _ := core.GetAWSKey(destination)
		Expect(key).To(Equal("/OBJECT/HERE.mp4"))
	})
})
