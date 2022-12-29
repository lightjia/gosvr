package core

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

var (
	entityMgr     *CEntityMgr
	entityMgrOnce sync.Once
)

type FuncEntityCreate func(context.Context, map[string]interface{}) Entity
type FuncComponentCreate func(Entity, map[string]interface{}) Component

type CEntityMgr struct {
	entityCretors    map[string]FuncEntityCreate
	entityComponents map[string][]string
	componentCretors map[string]FuncComponentCreate
	entitysLock      *sync.RWMutex
	entitys          map[int64]Entity
}

func (this *CEntityMgr) New() {
	this.entityCretors = make(map[string]FuncEntityCreate)
	this.entityComponents = make(map[string][]string)
	this.componentCretors = make(map[string]FuncComponentCreate)
	this.entitysLock = &sync.RWMutex{}
	this.entitys = make(map[int64]Entity)
}

func (this *CEntityMgr) RegisterEntity(name string, f FuncEntityCreate, comps []string) {
	this.entityCretors[name] = f
	this.entityComponents[name] = append(this.entityComponents[name], comps...)
}

func (this *CEntityMgr) RegisterComponent(name string, f FuncComponentCreate) {
	this.componentCretors[name] = f
}

func (this *CEntityMgr) NewEntity(ctx context.Context, name string, args map[string]interface{}) Entity {
	if v, ok := this.entityCretors[name]; ok {
		ret := v(ctx, args)
		ret.Context().SetValue(Entity_Kind, name)
		comps, _ := this.entityComponents[name]
		for _, comp := range comps {
			if c, cok := this.componentCretors[comp]; cok {
				ret.AddComp(c(ret, args))
			} else {
				LOG.Warn(fmt.Sprintf("new_entity_error entity_kind=%s,,comp_name=%s", name, comp))
			}
		}

		ret.Init(args)
		return ret
	}

	LOG.Warn(fmt.Sprintf("no_entity_creator kind=%s", name))
	return nil
}

func (this *CEntityMgr) AddEntity(ent Entity) {
	this.entitysLock.Lock()
	defer this.entitysLock.Unlock()
	entID, _ := ent.Context().Value(Entity_Id).(int64)
	this.entitys[entID] = ent
}

func (this *CEntityMgr) DelEntity(ent Entity) {
	this.entitysLock.Lock()
	defer this.entitysLock.Unlock()
	entID, _ := ent.Context().Value(Entity_Id).(int64)
	delete(this.entitys, entID)
}

func (this *CEntityMgr) GetEntityByID(entID int64) Entity {
	this.entitysLock.RLock()
	defer this.entitysLock.RUnlock()
	if v, ok := this.entitys[entID]; ok {
		return v
	}

	return nil
}

func (this *CEntityMgr) GetEntityByKind(kind string) (ret []Entity) {
	this.entitysLock.RLock()
	defer this.entitysLock.RUnlock()
	for _, v := range this.entitys {
		k, _ := v.Context().Value(Entity_Kind).(string)
		if k == kind {
			ret = append(ret, v)
		}
	}

	return ret
}

func (this *CEntityMgr) String() string {
	var sb strings.Builder
	sb.WriteString("[Entitys:")
	for k := range this.entityComponents {
		sb.WriteString(fmt.Sprintf(" %s", k))
	}
	sb.WriteString("]\n[Components:")
	for k := range this.componentCretors {
		sb.WriteString(fmt.Sprintf(" %s", k))
	}
	sb.WriteString("]")
	return sb.String()
}

func EntityMgr() *CEntityMgr {
	entityMgrOnce.Do(func() {
		entityMgr = &CEntityMgr{}
		entityMgr.New()
	})

	return entityMgr
}
