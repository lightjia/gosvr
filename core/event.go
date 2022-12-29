package core

import "reflect"

type EventCallback func(interface{})

type EventCreateFunc func(string) *CEvent

func CreateEvent(name string) *CEvent {
	return &CEvent{async: false, name: name}
}

func CreateAsyncEvent(name string) *CEvent {
	return &CEvent{async: true, name: name}
}

type CEvent struct {
	async bool
	name  string
	cbs   []EventCallback
}

func (this *CEvent) Dispatch(args interface{}) {
	for _, v := range this.cbs {
		if this.async {
			go v(args)
			continue
		}

		v(args)
	}
}

func (this *CEvent) RegisterEvent(cb EventCallback) {
	for i := 0; i < len(this.cbs); i++ {
		sf1 := reflect.ValueOf(this.cbs[i])
		sf2 := reflect.ValueOf(cb)
		if sf1.Pointer() == sf2.Pointer() {
			return
		}
	}

	this.cbs = append(this.cbs, cb)
}

func (this *CEvent) UnRegisterEvent(cb EventCallback) {
	for i := 0; i < len(this.cbs); i++ {
		sf1 := reflect.ValueOf(this.cbs[i])
		sf2 := reflect.ValueOf(cb)
		if sf1.Pointer() == sf2.Pointer() {
			if i == len(this.cbs)-1 {
				l := i - 1
				if l < 0 {
					l = 0
				}
				this.cbs = this.cbs[:l]
			} else {
				this.cbs = append(this.cbs[:i], this.cbs[i+1:]...)
			}

			break
		}
	}
}
