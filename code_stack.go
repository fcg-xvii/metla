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
	for i := s.Len() - 1; i >= 0; i-- {
		//res[i] = s.Pop()
		val := s.Pop()
		switch val.(type) {
		case splitter:
		default:
			res = append(res, val)
		}
	}
	s.commandCount = 0
	return res
}
