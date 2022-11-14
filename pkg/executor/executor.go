package executor

import (
	"github.com/utkarsh-pro/notion-gister/pkg/executor/cron"
	"github.com/utkarsh-pro/notion-gister/pkg/executor/none"
)

type Executor interface {
	AddFunc(string, func()) error
	Start()
}

func New(typ string) Executor {
	switch typ {
	case "cron":
		return cron.New()
	case "none":
		return none.New()
	default:
		return none.New()
	}
}
