package downloaders

import (
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
)

// DownloadFunc is a function type for the multiple
// possible ways to download the source file
type DownloadFunc func(logger lager.Logger, config gonfig.Gonfig, dbInstance db.Storage, jobID string) error

// GetDownloadFunc returns the download function
// based on the job source.
func GetDownloadFunc(jobSource string) DownloadFunc {
	if strings.Contains(jobSource, "amazonaws") {
		return S3Download
	} else if strings.HasPrefix(jobSource, "ftp://") {
		return FTPDownload
	}

	return HTTPDownload
}
