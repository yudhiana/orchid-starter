package directive

import (
	"orchid-starter/internal/bootstrap"
)

type Directive struct {
	DI *bootstrap.DirectInjection
}

func NewDirective(di *bootstrap.DirectInjection) *Directive {
	return &Directive{
		DI: di,
	}
}
