package prod

import (
	"fmt"
	"time"
)

func layoutFromMap(src map[string]interface{}) *storageLayout {
	res := &storageLayout{
		list: make([]*variable, 0, len(src)),
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

func (s *storageLayout) appendVariable(key string, val interface{}) (*variable, error) {
	//s.list = append(s.list, v)
	if _, check := s.findVariable(key); check {
		return nil, fmt.Errorf("Set variable error :: variable [%v] exists on current layout", key)
	} else {
		res := &variable{key, val}
		s.list = append(s.list, res)
		return res, nil
	}
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
		execStart: time.Now(),
		layouts:   []*storageLayout{layout},
		layout:    layout,
	}
}

type storage struct {
	execStart time.Time
	layouts   []*storageLayout
	layout    *storageLayout
}

func (s *storage) checkTimeout() bool {
	return time.Now().Sub(s.execStart) > time.Second*30
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

func (s *storage) appendVariable(name string, value interface{}) (*variable, error) {
	return s.layout.appendVariable(name, value)
}

/*func (s *storage) updateVariable(v *variable) (res *variable) {
	var check bool
	if res, check = s.findVariable(v.key); check {
		res.value = v.value
	} else {
		res = v
		s.layout.appendVariable(v)
	}
	return
}*/

type variable struct {
	key   string
	value interface{}
}
