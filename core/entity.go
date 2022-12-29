package core

import (
	"context"
	"fmt"
	"sync"
)

type Entity interface {
	New(context.Context)
	Init(map[string]interface{}) bool
	RegisterEvent(string, EventCallback, EventCreateFunc)
	UnRegisterEvent(string, EventCallback)
	DispatchEvent(string, interface{})
	Destroy()
	AddComp(Component)
	Comp(string) Component
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
	Context() *Context
}

type CEntityBase struct {
	log *Logger
	ctx *Context
}

func (this *CEntityBase) Init(args map[string]interface{}) bool {
	comps, _ := this.ctx.Value(Entity_Component).(map[string]Component)
	for _, c := range comps {
		c.Init(this, args)
	}

	return true
}

func (this *CEntityBase) Destroy() {
	comps, _ := this.ctx.Value(Entity_Component).(map[string]Component)
	compsLock, _ := this.ctx.Value(Entity_Component_Lock).(*sync.RWMutex)
	compsLock.RLock()
	defer compsLock.RUnlock()
	for _, c := range comps {
		c.Destroy()
	}
	EntityMgr().DelEntity(this)
}

func (this *CEntityBase) Comp(name string) Component {
	comps, _ := this.ctx.Value(Entity_Component).(map[string]Component)
	compsLock, _ := this.ctx.Value(Entity_Component_Lock).(*sync.RWMutex)
	compsLock.RLock()
	defer compsLock.RUnlock()
	if v, ok := comps[name]; ok {
		return v
	}
	return nil
}

func (this *CEntityBase) AddComp(c Component) {
	if c.Name() != "" {
		comps, _ := this.ctx.Value(Entity_Component).(map[string]Component)
		compsLock, _ := this.ctx.Value(Entity_Component_Lock).(*sync.RWMutex)
		compsLock.Lock()
		defer compsLock.Unlock()
		comps[c.Name()] = c
	}
}

func (this *CEntityBase) New(ctx context.Context) {
	if !IsLightContext(ctx) {
		ctx = NewContext(ctx)
	}

	this.ctx, _ = ctx.(*Context)
	this.ctx.SetValue(Entity_Id, SnowID().NextVal())
	this.ctx.SetValue(Entity_Component, make(map[string]Component))
	this.ctx.SetValue(Entity_Event, make(map[string]*CEvent))
	this.ctx.SetValue(Entity_Event_Lock, &sync.RWMutex{})
	this.ctx.SetValue(Entity_Component_Lock, &sync.RWMutex{})
	this.log = &LOG
	EntityMgr().AddEntity(this)
}

func (this *CEntityBase) DispatchEvent(evName string, evParam interface{}) {
	events, _ := this.ctx.Value(Entity_Event).(map[string]*CEvent)
	eventsLock, _ := this.ctx.Value(Entity_Event_Lock).(*sync.RWMutex)
	eventsLock.RLock()
	defer eventsLock.RUnlock()
	if ev, ok := events[evName]; ok {
		ev.Dispatch(evParam)
	}
}

func (this *CEntityBase) RegisterEvent(evName string, evCb EventCallback, create EventCreateFunc) {
	events, _ := this.ctx.Value(Entity_Event).(map[string]*CEvent)
	eventsLock, _ := this.ctx.Value(Entity_Event_Lock).(*sync.RWMutex)
	eventsLock.Lock()
	defer eventsLock.Unlock()
	if ev, ok := events[evName]; ok {
		ev.RegisterEvent(evCb)
	} else {
		if create != nil {
			ev := create(evName)
			ev.RegisterEvent(evCb)
			events[evName] = ev
		}
	}
}

func (this *CEntityBase) UnRegisterEvent(evName string, evCb EventCallback) {
	events, _ := this.ctx.Value(Entity_Event).(map[string]*CEvent)
	eventsLock, _ := this.ctx.Value(Entity_Event_Lock).(*sync.RWMutex)
	eventsLock.Lock()
	defer eventsLock.Unlock()
	if ev, ok := events[evName]; ok {
		ev.UnRegisterEvent(evCb)
	}
}

func (this *CEntityBase) Context() *Context {
	return this.ctx
}

func (this *CEntityBase) Debug(format string, v ...interface{}) {
	this.log.Debug(fmt.Sprintf("[%s-%d]%s", this.ctx.Value(Entity_Kind), this.ctx.Value(Entity_Id), fmt.Sprintf(format, v...)))
}

func (this *CEntityBase) Info(format string, v ...interface{}) {
	this.log.Info(fmt.Sprintf("[%s-%d]%s", this.ctx.Value(Entity_Kind), this.ctx.Value(Entity_Id), fmt.Sprintf(format, v...)))
}

func (this *CEntityBase) Error(format string, v ...interface{}) {
	this.log.Error(fmt.Sprintf("[%s-%d]%s", this.ctx.Value(Entity_Kind), this.ctx.Value(Entity_Id), fmt.Sprintf(format, v...)))
}
