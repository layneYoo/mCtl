package marathon

import (
	"github.com/layneYoo/mCtl/check"
)

type Selector interface {
	Select(args []string)
}

type Category struct {
	Actions map[string]Action
}

func (c Category) Select(args []string) {
	check.Check(len(args) > 0, "must specify sub-action")
	if action, ok := c.Actions[args[0]]; !ok {
		check.Usage()
	} else {
		action.Apply(args[1:])
	}
}
