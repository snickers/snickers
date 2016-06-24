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
	out, _ := os.Create(d.tempPath + "output.mp4")
	defer out.Close()

	resp, _ := http.Get(d.job.Source)
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
}
