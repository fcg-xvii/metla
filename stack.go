package metla

import stack "github.com/golang-collections/collections/stack"

type stackStorage struct {
	names  []string
	values []interface{}
	marks  *stack.Stack
}

func (s *stackStorage) PushVariable(name string, value interface{}) (result int) {
	result = len(s.names)
	s.names = append(s.names, name)
	s.values = append(s.values, value)
	return
}

func (s *stackStorage) PushMark() {
	s.marks.Push(len(s.names))
}

func (s *stackStorage) PopMark() {
	if s.marks.Len() == 0 {
		return
	}
	endIndex := s.marks.Pop().(int)
	if endIndex < len(s.names) {
		s.names = s.names[:endIndex]
		s.values = s.values[:endIndex]
	}
}

func (s *stackStorage) FindName(name string) (index int, check bool) {
	for i := len(s.names) - 1; i >= 0; i-- {
		if s.names[i] == name {
			return i, true
		}
	}
	return
}
