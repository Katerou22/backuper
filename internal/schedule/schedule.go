package schedule

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

type Task struct {
	task     func()
	cronTime string
}

type Schedule struct {
	tasks []Task
	cron  *cron.Cron
}

func NewScheduler() *Schedule {

	var schedule Schedule

	c := cron.New(cron.WithSeconds())

	schedule.cron = c

	return &schedule

}

func (schedule *Schedule) AddTask(task func(), cronTime string) {

	fmt.Println("Added Task on " + cronTime)
	schedule.cron.AddFunc(cronTime, func() {
		fmt.Printf("Running task at %s\n", time.Now())
		task()
	})
}

func (schedule *Schedule) Run() {
	fmt.Println("Running schedule at " + time.Now().String())
	schedule.cron.Start()
}
