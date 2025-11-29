package singleton

import (
	"reflect"
	"sync"
)

type Struct = interface{}

var lock sync.Mutex
var instances = make(map[reflect.Type]interface{})

func GetInstance[T Struct](create ...func() *T) *T {
	lock.Lock()
	defer lock.Unlock()

	if len(create) == 0 {
		create = append(create, func() *T { return new(T) })
	} else if len(create) > 1 {
		panic("Multiple parameters")
	}

	t_type := reflect.TypeFor[T]()
	if t_type.Kind() != reflect.Struct {
		panic("T(type T) does not satisfy struct")
	}

	if instance, ok := instances[t_type]; ok {
		return instance.(*T)
	}

	instances[t_type] = create[0]()

	return instances[t_type].(*T)
}
