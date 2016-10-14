package uploaders

import (
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/snickers/snickers/db"
)

// UploadFunc is a function type for the multiple
// possible ways to upload the source file
type UploadFunc func(logger lager.Logger, dbInstance db.Storage, jobID string) error

// GetUploadFunc returns the upload function
// based on the job source.
func GetUploadFunc(jobDestination string) UploadFunc {
	if strings.HasPrefix(jobDestination, "ftp://") {
		return FTPUpload
	}

	return S3Upload
}
