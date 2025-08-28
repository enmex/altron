package scheduler

import "time"

type Scheduler struct {
	ticker  *time.Ticker
	handler func() error
	err     chan error
}

func NewScheduler(interval time.Duration, handler func() error) *Scheduler {
	return &Scheduler{
		ticker:  time.NewTicker(interval),
		handler: handler,
		err:     make(chan error),
	}
}

func (s *Scheduler) Start() {
	go func() {
		for {
			<-s.ticker.C
			if err := s.handler(); err != nil {
				s.err <- err
			}
		}
	}()
}

func (s *Scheduler) Error() <-chan error {
	return s.err
}
