package downloaders

import (
	"os"

	"code.cloudfoundry.org/lager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/helpers"
)

// S3Download downloads the file from S3 bucket. Job Source should be
// in format: http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT
func S3Download(logger lager.Logger, config gonfig.Gonfig, dbInstance db.Storage, jobID string) error {
	log := logger.Session("s3-download")
	log.Info("start", lager.Data{"job": jobID})
	defer log.Info("finished")

	job, err := SetupJob(jobID, dbInstance, config)
	if err != nil {
		log.Error("setting-up-job", err)
		return err
	}

	file, err := os.Create(job.LocalDestination)
	if err != nil {
		return err
	}

	err = helpers.SetAWSCredentials(job.Source)
	if err != nil {
		return err
	}

	bucket, err := helpers.GetAWSBucket(job.Source)
	if err != nil {
		return err
	}

	key, err := helpers.GetAWSKey(job.Source)
	if err != nil {
		return err
	}

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	objInput := s3.GetObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)}

	_, err = downloader.Download(file, &objInput)

	return err
}
