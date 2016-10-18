package encoders

import (
	"code.cloudfoundry.org/lager"
	"github.com/snickers/hls/segmenter"
	"github.com/snickers/snickers/db"
)

// HLSEncode function is responsible for encoding adaptive bitrate outputs
func HLSEncode(logger lager.Logger, dbInstance db.Storage, jobID string) error {
	log := logger.Session("hls-encode")
	log.Info("started", lager.Data{"job": jobID})
	defer log.Info("finished")
	var cfg segmenter.HLSConfig
	if cfg.SourceFile == "" {
		return nil
	}
	return nil
}
