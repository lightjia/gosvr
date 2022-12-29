package main

import (
	"context"
	"fmt"
	"light/component"
	"light/core"
	"light/entity"
)

func main() {
	fmt.Println(core.EntityMgr())
	a := core.EntityMgr().NewEntity(context.Background(), "A", nil)
	if b, ok := a.(*entity.CAEntity); ok {
		fmt.Println("hello lightjia", b.Context())
		b.DispatchEvent("Test", nil)
		c := a.Comp("A")
		d, _ := c.(*component.CAComponent)
		el := len(core.EntityMgr().GetEntityByKind("A"))
		fmt.Println("hello component", d.Name(), el)
	} else {
		fmt.Printf("a=%v", a)
	}

	core.LOG.Debug("hello world")
}
