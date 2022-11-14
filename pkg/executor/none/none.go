package none

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type None struct {
	fns []func()
}

func New() *None {
	return &None{}
}

func (n *None) AddFunc(spec string, fn func()) error {
	n.fns = append(n.fns, fn)
	return nil
}

func (n *None) Start() {
	logrus.Info("Starting none executor")
	var wg sync.WaitGroup
	for _, fn := range n.fns {
		wg.Add(1)

		go func(fn func()) {
			fn()
			wg.Done()
		}(fn)
	}

	wg.Wait()
}
