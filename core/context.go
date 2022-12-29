package core

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type Context struct {
	tagsLock *sync.RWMutex
	tags     map[interface{}]interface{}
	context.Context
}

func NewContext(ctx context.Context) *Context {
	tagsLock := &sync.RWMutex{}
	ctx = context.WithValue(ctx, ContextTagsLock, tagsLock)
	return &Context{
		tagsLock: tagsLock,
		Context:  ctx,
		tags:     map[interface{}]interface{}{isLightContext: true},
	}
}

func (c *Context) Lock() {
	c.tagsLock.Lock()
}

func (c *Context) Unlock() {
	c.tagsLock.Unlock()
}

func (*Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*Context) Done() <-chan struct{} {
	return nil
}

func (*Context) Err() error {
	return nil
}

func (c *Context) Value(key interface{}) interface{} {
	c.tagsLock.RLock()
	defer c.tagsLock.RUnlock()
	if c.tags == nil {
		c.tags = make(map[interface{}]interface{})
	}

	if v, ok := c.tags[key]; ok {
		return v
	}
	return c.Context.Value(key)
}

func (c *Context) SetValue(key, val interface{}) {
	c.tagsLock.Lock()
	defer c.tagsLock.Unlock()

	if c.tags == nil {
		c.tags = make(map[interface{}]interface{})
	}
	c.tags[key] = val
}

// DeleteKey delete the kv pair by key.
func (c *Context) DeleteKey(key interface{}) {
	c.tagsLock.Lock()
	defer c.tagsLock.Unlock()

	if c.tags == nil || key == nil {
		return
	}
	delete(c.tags, key)
}

func (c *Context) String() string {
	return fmt.Sprintf("%v.WithValue(%v)", c.Context, c.tags)
}

func (c *Context) RLock() {
	c.tagsLock.RLock()
}

func (c *Context) RUnlock() {
	c.tagsLock.RUnlock()
}

func WithValue(parent context.Context, key, val interface{}) *Context {
	if key == nil {
		panic("nil key")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	tags := make(map[interface{}]interface{})
	tags[key] = val
	return &Context{Context: parent, tags: tags, tagsLock: &sync.RWMutex{}}
}

func WithLocalValue(ctx *Context, key, val interface{}) *Context {
	if key == nil {
		panic("nil key")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	if ctx.tags == nil {
		ctx.tags = make(map[interface{}]interface{})
	}

	ctx.tags[key] = val
	return ctx
}

// IsLightContext checks whether a context is core.Context.
func IsLightContext(ctx context.Context) bool {
	ok := ctx.Value(isLightContext)
	return ok != nil
}
