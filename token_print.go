package metla

import (
	"fmt"
	"io"
)

type tokenPrint struct {
	val token
}

func (s *tokenPrint) Val() interface{} {
	return s.val
}

func (s *tokenPrint) Data(w io.Writer, sto *storage) error {
	fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&", s.val)
	return s.val.Data(w, sto)
}

func (s *tokenPrint) String() string {
	return "[tokenPrint :: { " + s.String() + " }]"
}

func (s *tokenPrint) IsExecutable() bool { return false }
