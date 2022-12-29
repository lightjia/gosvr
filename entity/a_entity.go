package entity

import (
	"context"
	"light/core"
)

type CAEntity struct {
	*core.CEntityBase
}

func (this *CAEntity) OnEvent(evArg interface{}) {
	this.Debug("Hello in onevent")
}

func newA(ctx context.Context, args map[string]interface{}) core.Entity {
	ret := &CAEntity{&core.CEntityBase{}}
	ret.New(ctx)
	ret.RegisterEvent("Test", ret.OnEvent, core.CreateEvent)
	return ret
}

func init() {
	core.EntityMgr().RegisterEntity("A", newA, []string{"A", "B"})
}
