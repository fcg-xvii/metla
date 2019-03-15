package metla

import "github.com/golang-collections/collections/stack"

func newCodeStack() *codeStack {
	return &codeStack{stack.New(), 0}
}

type codeStack struct {
	*stack.Stack
	commandCount int
}

func (s *codeStack) Push(value interface{}) {
	if _, check := value.(*execCommand); check {
		s.commandCount++
	}
	s.Stack.Push(value)
}

func (s *codeStack) Pop() interface{} {
	if _, check := s.Peek().(execCommand); check {
		s.commandCount--
	}
	return s.Stack.Pop()
}

func (s *codeStack) Flush() []interface{} {
	res := make([]interface{}, 0, s.Len())
	for s.Len() > 0 {
		res = append(res, s.Pop())
	}
	s.commandCount = 0
	return res
}
