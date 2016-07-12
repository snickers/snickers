package lib

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/flavioribeiro/snickers/db"
)

// S3Upload sends the file to S3 bucket. Job Destination should be
// in format: http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT
func S3Upload(jobID string) error {
	dbInstance, err := db.GetDatabase()
	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return err
	}

	file, err := os.Open(job.LocalSource)
	if err != nil {
		return err
	}

	err = SetAWSCredentials(job.Destination)
	if err != nil {
		return err
	}

	bucket, err := GetAWSBucket(job.Destination)
	if err != nil {
		return err
	}

	key, err := GetAWSKey(job.Destination)
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	_, err = uploader.Upload(&s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	return nil
}
