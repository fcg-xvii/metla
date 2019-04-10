package prod

import (
	"fmt"

	_ "github.com/fcg-xvii/containers"
)

type executer interface {
	Exec(*tplExec) error
}

type getter interface {
	Get(*tplExec) error
}

type setter interface {
	Set(*tplExec, interface{}) error
}

type execText struct {
	src []byte
}

func (s execText) Exec(exec *tplExec) error {
	return exec.Write(s.src)
}

type execBase struct {
	line, pos int
}

func (s *execBase) Exec() (err error) {
	return fmt.Errorf("ALLLLLLLLLLLLLLLLLLLL")
}
