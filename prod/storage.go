package prod

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

/////////////////////////////////////////////

type execStorage struct {
	values []interface{}
	store  *storage
}
