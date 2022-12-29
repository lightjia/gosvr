package entity

import (
	"context"
	"light/core"
)

type CGlobalEntity struct {
	*core.CEntityBase
}

func (this *CGlobalEntity) OnEvent(evArg interface{}) {
	this.Debug("Hello in onevent")
}

func newGlobal(ctx context.Context, args map[string]interface{}) core.Entity {
	ret := &CGlobalEntity{&core.CEntityBase{}}
	ret.New(ctx)
	ret.RegisterEvent("Test", ret.OnEvent, core.CreateEvent)
	return ret
}

func init() {
	core.EntityMgr().RegisterEntity("Global", newGlobal, []string{"A", "B"})
}
