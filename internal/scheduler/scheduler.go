package scheduler

import (
	"awesomeProject3/internal/config"
	"log"
	"time"
)

type Job func()

type Scheduler struct {
	interval time.Duration
	job      Job
	stop     chan struct{}
}

func New(job Job) *Scheduler {
	return &Scheduler{
		interval: config.ScrapeInterval,
		job:      job,
		stop:     make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	log.Printf("Scheduler started, interval: %s", s.interval)
	go s.runOnce()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			go s.runOnce()
		}
	}
}

func (s *Scheduler) runOnce() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Scheduler job panic: %v", r)
		}
	}()
	s.job()
}

func (s *Scheduler) Stop() {
	close(s.stop)
}
