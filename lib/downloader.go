package lib

import (
	"github.com/flavioribeiro/snickers/types"
	"io"
	"net/http"
	"os"
)

type Downloader struct {
	tempPath string
	job      types.Job
}

func (d *Downloader) Start() {
	out, err := os.Create(d.tempPath + "output.mp4")
	defer out.Close()

	resp, err := http.Get(job.Source)
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
}
