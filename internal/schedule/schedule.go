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
		fmt.Println("Running task at", time.Now())
		go task()
	})
}

func (schedule *Schedule) Run() {
	fmt.Println("Running schedule")
	schedule.cron.Start()
}
