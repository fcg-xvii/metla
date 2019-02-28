package metla

import (
	"fmt"
	"reflect"
)

func layoutFromMap(src map[string]interface{}) *storageLayout {
	res := &storageLayout{
		list: make([]*variable, len(src)),
	}
	count := 0
	for key, val := range src {
		res.list = append(res.list, &variable{key, val})
		count++
	}
	return res
}

type storageLayout struct {
	list []*variable
}

func (s *storageLayout) appendVariable(key string, val interface{}) (res *variable, err error) {
	if _, check := s.findVariable(key); check {
		err = fmt.Errorf("Set variable error :: variable [%v] exists on current layout", key)
	} else {
		res = &variable{key, val}
		s.list = append(s.list, res)
	}
	return
}

func (s *storageLayout) findVariable(key string) (res *variable, check bool) {
	for i := len(s.list) - 1; i >= 0; i-- {
		res = s.list[i]
		if res.key == key {
			return res, true
		}
	}
	return
}

//////////////////////////////////////////////////////////////////////

func newStorage(src map[string]interface{}) *storage {
	layout := layoutFromMap(src)
	return &storage{
		layouts: []*storageLayout{layout},
		layout:  layout,
	}
}

type storage struct {
	layouts []*storageLayout
	layout  *storageLayout
}

func (s *storage) newLayout() {
	layout := new(storageLayout)
	s.layouts = append(s.layouts, layout)
	s.layout = layout
}

func (s *storage) dropLayout() {
	s.layouts = s.layouts[:len(s.layouts)-1]
	s.layout = s.layouts[len(s.layouts)-1]
}

func (s *storage) findVariable(key string) (res *variable, check bool) {
	if res, check = s.layout.findVariable(key); !check && len(s.layouts) > 1 {
		for i := len(s.layouts) - 2; i >= 0; i-- {
			if res, check = s.layouts[i].findVariable(key); check {
				return
			}
		}
	}
	return
}

func (s *storage) appendValue(key string, value interface{}) (res *variable, err error) {
	return s.layout.appendVariable(key, value)
}

func (s *storage) setValue(key string, value interface{}) {

}

type variable struct {
	key   string
	value interface{}
}

func (s *variable) Kind() reflect.Kind {
	return reflect.ValueOf(s.value).Kind()
}

func (s *variable) IsNil() bool {
	return reflect.ValueOf(s.value).IsNil()
}
