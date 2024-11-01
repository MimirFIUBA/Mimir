package triggers

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	gocronScheduler gocron.Scheduler
}

func NewScheduler() (*Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	return &Scheduler{s}, nil
}

func (s *Scheduler) ScheduleTrigger(cronExpresion string, trigger Trigger) {
	s.gocronScheduler.NewJob(
		gocron.CronJob(cronExpresion, false),
		gocron.NewTask(
			func() {
				trigger.Update(Event{
					SenderId:  "Scheduler",
					Type:      SCHEDULER_ACTIVE,
					Timestamp: time.Now(),
				})
			},
		),
	)
}

func (s *Scheduler) Start() {
	s.gocronScheduler.Start()
}

func (s *Scheduler) Shutdown() {
	s.gocronScheduler.Shutdown()
}
