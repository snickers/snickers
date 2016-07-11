package lib

import (
	"log"
	"net/url"
	"os"
	"strings"

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
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	return nil
}

// GetAWSKey grabs the path and filename for destination
func GetAWSKey(jobSource string) (string, error) {
	parsedURL, err := url.Parse(jobSource)
	if err != nil {
		return "", err
	}

	return parsedURL.Path, nil
}

// GetAWSBucket grabs the bucket from a given s3 source
func GetAWSBucket(jobSource string) (string, error) {
	parsedURL, err := url.Parse(jobSource)
	if err != nil {
		return "", err
	}

	region := strings.Split(parsedURL.Host, ".")[0]
	return region, nil
}

// SetAWSCredentials will parse the job source and set the credentials
// on environment variables
func SetAWSCredentials(jobSource string) error {
	parsedURL, err := url.Parse(jobSource)
	if err != nil {
		return err
	}
	os.Setenv("AWS_ACCESS_KEY_ID", parsedURL.User.Username())
	password, _ := parsedURL.User.Password()
	os.Setenv("AWS_SECRET_ACCESS_KEY", password)
	return nil
}
