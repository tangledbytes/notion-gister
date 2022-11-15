package cron

import (
	rcron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/utkarsh-pro/notion-gister/pkg/config"
)

type Cron struct {
	cron *rcron.Cron
}

func New() *Cron {
	return &Cron{
		cron: rcron.New(
			rcron.WithLocation(config.Timezone()),
		),
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
