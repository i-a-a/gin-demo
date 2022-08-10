package pkg

import (
	cron "github.com/robfig/cron/v3"
)

var (
	Cron = &myCron{}
)

type myCron struct {
	*cron.Cron
}

// spce 可以是"* * * * *", 也可以是 "@every 10s" 这种
func (c *myCron) AddFunc(spec string, fn func()) {
	if c.Cron == nil {
		c.Cron = cron.New()
	}
	c.Cron.AddFunc(spec, fn)
}

// 在确认所有调用 AddFunc 完成后再Start
func (c *myCron) Start() {
	c.Cron.Start()
}
