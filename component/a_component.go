package component

import "light/core"

type CAComponent struct {
	*core.CComponentBase
}

const (
	name = "A"
)

func (this *CAComponent) Name() string {
	return name
}

func newA(entity core.Entity, args map[string]interface{}) core.Component {
	ret := &CAComponent{&core.CComponentBase{}}
	return ret
}

func init() {
	core.EntityMgr().RegisterComponent(name, newA)
}
