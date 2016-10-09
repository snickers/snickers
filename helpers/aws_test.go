package helpers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Helpers", func() {

	It("Should get bucket from URL Destination", func() {
		destination := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT"
		bucket, _ := GetAWSBucket(destination)
		Expect(bucket).To(Equal("BUCKET"))
	})

	It("Should set credentials from URL Destination", func() {
		destination := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT"
		SetAWSCredentials(destination)
		Expect(os.Getenv("AWS_ACCESS_KEY_ID")).To(Equal("AWSKEY"))
		Expect(os.Getenv("AWS_SECRET_ACCESS_KEY")).To(Equal("AWSSECRET"))
	})

	It("Should get path and filename from URL Destination", func() {
		destination := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT/HERE.mp4"
		key, _ := GetAWSKey(destination)
		Expect(key).To(Equal("/OBJECT/HERE.mp4"))
	})
})
