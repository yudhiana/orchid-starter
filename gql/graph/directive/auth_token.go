package directive

import (
	"orchid-starter/internal/bootstrap/container"
)

type Directive struct {
	DI *container.DirectInjection
}

func NewDirective(di *container.DirectInjection) *Directive {
	return &Directive{
		DI: di,
	}
}
