package cron

import (
	rcron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Cron struct {
	cron *rcron.Cron
}

func New() *Cron {
	return &Cron{
		cron: rcron.New(),
	}
}

func (c *Cron) AddFunc(spec string, fn func()) error {
	_, err := c.cron.AddFunc(spec, fn)
	return err
}

func (c *Cron) Start() {
	logrus.Info("Starting cron executor")
	c.cron.Run()
}
