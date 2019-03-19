package metla

import (
	"fmt"
	"reflect"
)

func layoutFromMap(src map[string]interface{}) *storageLayout {
	res := &storageLayout{
		list: make([]*variable, 0, len(src)),
	}
	count := 0
	for key, val := range src {
		res.list = append(res.list, &variable{key, val, true})
		count++
	}
	return res
}

type storageLayout struct {
	list []*variable
}

func (s *storageLayout) appendVariable(v *variable) {
	s.list = append(s.list, v)
	/*if _, check := s.findVariable(key); check {
		err = fmt.Errorf("Set variable error :: variable [%v] exists on current layout", key)
	} else {
		res = &variable{key, val}
		s.list = append(s.list, res)
	}
	return*/
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
	//fmt.Println("MAP.........", src)
	layout := layoutFromMap(src)
	//fmt.Println("LAY.........", layout)
	return &storage{
		layouts: []*storageLayout{layout},
		layout:  layout,
	}
}

type storage struct {
	layouts []*storageLayout
	layout  *storageLayout
}

func (s *storage) newLayout() *storageLayout {
	layout := new(storageLayout)
	s.layouts = append(s.layouts, layout)
	s.layout = layout
	return layout
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

func (s *storage) appendVariable(v *variable) {
	s.layout.appendVariable(v)
}

func (s *storage) updateVariable(v *variable) (res *variable) {

	var check bool
	if res, check = s.findVariable(v.key); check {
		res.value = v.value
	} else {
		res = v
		s.layout.appendVariable(v)
	}
	return
}

type variable struct {
	key    string
	value  interface{}
	stored bool
}

func (s *variable) Kind() reflect.Kind {
	return reflect.ValueOf(s.value).Kind()
}

func (s *variable) IsNil() bool {
	return reflect.ValueOf(s.value).IsNil()
}

func (s *variable) String() string { return fmt.Sprint(s.value) }
