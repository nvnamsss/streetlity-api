package event

import (
	"fmt"
	"reflect"
	"runtime"
)

type Event struct {
	Async     bool
	callbacks []interface{}
}

func (e *Event) Subscribe(callback interface{}) {
	e.callbacks = append(e.callbacks, callback)
	v := reflect.ValueOf(e.callbacks[0])

	fmt.Println(v.Kind() == reflect.Func)
}

func (e *Event) Unsubscribe(callback interface{}) {
	fmt.Println(runtime.FuncForPC(reflect.ValueOf(callback).Pointer()).Name())
}

func (e *Event) Invoke() {
	for _, item := range e.callbacks {
		if f, ok := item.(func()); ok {
			if e.Async {
				go f()
			} else {
				f()
			}
		}
	}

}

func NewEvent() *Event {
	e := new(Event)
	e.Async = true
	return e
}
