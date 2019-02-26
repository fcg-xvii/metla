package metla

import (
	"fmt"
	"io"
	"strconv"
)

type tokenText struct {
	src []byte
}

func (s *tokenText) Data(w io.Writer, sto *storage) (err error) {
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!", string(s.src))
	_, err = w.Write(s.src)
	return
}

func (s *tokenText) String() string {
	return "[tokenText :: { []byte :: len:" + strconv.Itoa(len(s.src)) + " }]"
}

func (s *tokenText) IsExecutable() bool { return false }
