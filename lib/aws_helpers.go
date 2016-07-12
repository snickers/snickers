package lib

import (
	"net/url"
	"os"
	"strings"
)

// GetAWSKey grabs the path and filename for destination
func GetAWSKey(jobSource string) (string, error) {
	parsedURL, err := url.Parse(jobSource)
	if err != nil {
		return "", err
	}

	return parsedURL.Path, nil
}

// GetAWSBucket grabs the bucket from a given s3 url
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
