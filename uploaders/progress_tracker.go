package uploaders

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

type ProgressTracker struct {
	current int64
	Total   int64

	job *types.Job
	db  db.Storage

	finishOnce sync.Once
	finish     chan struct{}
	isFinish   bool

	startValue   int64
	currentValue int64
}

func (pt *ProgressTracker) Start() *ProgressTracker {
	pt.startValue = pt.current
	pt.Update()
	go pt.refresher()
	return pt
}

func (pt *ProgressTracker) Finish() {
	pt.finishOnce.Do(func() {
		close(pt.finish)
		pt.isFinish = true
	})
}

func (pt *ProgressTracker) Update() {
	c := atomic.LoadInt64(&pt.current)
	if c != pt.currentValue {
		pt.currentValue = c
	}

	percent := (float64(c) / float64(pt.Total)) * float64(100)

	pt.job.Progress = fmt.Sprintf("%6.02f%%", percent)
	pt.db.UpdateJob(pt.job.ID, *pt.job)

	if c >= pt.Total && pt.isFinish != true {
		pt.Finish()
	}
}

func (pt *ProgressTracker) Get() int64 {
	c := atomic.LoadInt64(&pt.current)
	return c
}

func (pt *ProgressTracker) Increment() int {
	return pt.Add(1)
}

func (pt *ProgressTracker) Add(add int) int {
	return int(pt.Add64(int64(add)))
}

func (pt *ProgressTracker) Add64(add int64) int64 {
	return atomic.AddInt64(&pt.current, add)
}

func (pt *ProgressTracker) Set(current int) *ProgressTracker {
	return pt.Set64(int64(current))
}

func (pt *ProgressTracker) Set64(current int64) *ProgressTracker {
	atomic.StoreInt64(&pt.current, current)
	return pt
}

func (pt *ProgressTracker) Read(p []byte) (n int, err error) {
	n = len(p)
	pt.Add(n)
	return
}

func (pt *ProgressTracker) Write(p []byte) (n int, err error) {
	n = len(p)
	pt.Add(n)
	return
}

func NewProgressTracker(total int, job *types.Job, dbInstance db.Storage) *ProgressTracker {
	return NewProgressTracker64(int64(total), job, dbInstance)
}

func NewProgressTracker64(total int64, job *types.Job, dbInstance db.Storage) *ProgressTracker {
	pt := &ProgressTracker{
		Total:        total,
		job:          job,
		currentValue: -1,
		db:           dbInstance,
		finish:       make(chan struct{}),
	}
	return pt
}

func (pt *ProgressTracker) refresher() {
	for {
		select {
		case <-pt.finish:
			return
		case <-time.After(time.Millisecond * 200):
			pt.Update()
		}
	}
}
