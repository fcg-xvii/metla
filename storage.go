package metla

import (
	"errors"
	"fmt"
	_ "time"
)

type variable struct {
	layout int
	key    string
	global bool
}

func (s *variable) String() string {
	return fmt.Sprintf("{ key: %v, layout: %v, global: %v }", s.key, s.layout, s.global)
}

type storage struct {
	list   []*variable
	layout int
}

func (s *storage) findVariable(key string) int {
	layout := s.layout
	for layout >= 0 {
		if index := s.findVariableInLayout(key, layout); index >= 0 {
			return index
		}
		layout--
	}
	return -1
}

func (s *storage) findVariableInLayout(key string, layout int) int {
	for i, v := range s.list {
		if v.key == key && v.layout == layout {
			return i
		}
	}
	return -1
}

func (s *storage) initVariable(key string) int {
	if index := s.findVariable(key); index >= 0 {
		return index
	}
	s.list = append(s.list, &variable{0, key, true})
	return len(s.list) - 1
}

func (s *storage) setVariable(key string) (int, error) {
	if s.findVariableInLayout(key, s.layout) != -1 {
		return -1, errors.New("Variable already exists in current layout")
	}
	v := &variable{s.layout, key, false}
	s.list = append(s.list, v)
	return len(s.list) - 1, nil
}

func (s *storage) saveInEmptyIndex(v *variable) int {
	for i, v := range s.list {
		if v == nil {
			s.list[i] = v
			return i
		}
	}
	s.list = append(s.list, v)
	return len(s.list) - 1
}

func (s *storage) globalIndexes() (res []int) {
	for i, v := range s.list {
		if v != nil && v.global {
			res = append(res, i)
		}
	}
	return
}

func (s *storage) globalKeys() map[string]int {
	res := make(map[string]int)
	for i, v := range s.list {
		if v.global {
			res[v.key] = i
		}
	}
	return res
}

func (s *storage) incLayout() {
	s.layout++
}

func (s *storage) decLayout() {
	/*for i, v := range s.list {
		if v.layout == s.layout {
			s.list[i] = nil
		}
	}*/
	s.layout--
}

func (s *storage) execStorage(vals map[string]interface{}) *execStorage {
	res := &execStorage{
		values: make([]interface{}, len(s.list)),
		store:  s,
	}
	if global := s.globalIndexes(); len(global) > 0 {
		for key, val := range vals {
			for _, index := range global {
				if key == s.list[index].key {
					res.values[index] = val
				}
			}
		}
	}
	//fmt.Println("STO_EXEC", res)
	return res
}

/////////////////////////////////////////////

type execStorage struct {
	values []interface{}
	store  *storage
}

func (s *execStorage) setValue(index int, value interface{}) {
	s.values[index] = value
}

func (s *execStorage) getValue(index int) interface{} {
	return s.values[index]
}

func (s *execStorage) globalMapNotNil() map[string]interface{} {
	res := make(map[string]interface{})
	for _, v := range s.store.globalIndexes() {
		if val := s.values[v]; val != nil {
			res[s.store.list[v].key] = val
		}
	}
	return res
}

func (s *execStorage) compare(sto *execStorage) {
	sVals, globalKeys := sto.globalMapNotNil(), s.store.globalKeys()
	delete(sVals, "params")
	for key, val := range sVals {
		if index, check := globalKeys[key]; check {
			s.values[index] = val
		}
	}
}
