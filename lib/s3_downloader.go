package lib

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/flavioribeiro/snickers/db"
)

// S3Download downloads the file from S3 bucket. Job Source should be
// in format: http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT
func S3Download(jobID string) error {
	dbInstance, err := db.GetDatabase()
	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return err
	}

	file, err := os.Open(job.LocalDestination)
	if err != nil {
		return err
	}

	err = SetAWSCredentials(job.Source)
	if err != nil {
		return err
	}

	bucket, err := GetAWSBucket(job.Source)
	if err != nil {
		return err
	}

	key, err := GetAWSKey(job.Source)
	if err != nil {
		return err
	}

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	_, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		return err
	}
	return nil
}
