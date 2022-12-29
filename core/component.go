package core

type Component interface {
	Init(Entity, map[string]interface{}) bool
	Destroy()
	Name() string
}

type CComponentBase struct {
	entity Entity
}

func (this *CComponentBase) Init(entity Entity, args map[string]interface{}) bool {
	this.entity = entity
	return true
}

func (this *CComponentBase) Name() string {
	return ""
}

func (this *CComponentBase) Destroy() {
}
