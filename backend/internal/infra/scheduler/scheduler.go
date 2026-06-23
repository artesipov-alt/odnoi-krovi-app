package scheduler

import (
	"context"
	"time"
)

type Job interface {
	Run(ctx context.Context)
}

type Scheduler struct {
	jobs   []scheduledJob
	cancel context.CancelFunc
}

type scheduledJob struct {
	job      Job
	interval time.Duration
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Register(job Job, interval time.Duration) {
	s.jobs = append(s.jobs, scheduledJob{job, interval})
}

func (s *Scheduler) Start() {
	ctx := context.Background()
	ctx, s.cancel = context.WithCancel(ctx)
	for _, sj := range s.jobs {
		go s.runJob(ctx, sj)
	}
}

func (s *Scheduler) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Scheduler) runJob(ctx context.Context, sj scheduledJob) {
	ticker := time.NewTicker(sj.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sj.job.Run(ctx)
		}
	}
}
