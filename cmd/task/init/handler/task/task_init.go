package task

import (
	modelInit "orchid-starter/cmd/task/init/model"
	"orchid-starter/internal/bootstrap"

	"github.com/yudhiana/logos"
)

type initTask struct {
	di   *bootstrap.DirectInjection
	init modelInit.Init
	log  *logos.LogEntry
}

func NewTask(di *bootstrap.DirectInjection, init modelInit.Init) *initTask {
	return &initTask{
		di:   di,
		init: init,
		log:  logos.NewLogger(),
	}
}

func (base *initTask) Start() (err error) {
	base.log.Info("Initialize task init")
	return nil
}
