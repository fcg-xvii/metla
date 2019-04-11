package prod

import (
	"fmt"

	_ "github.com/fcg-xvii/containers"
)

type coordinator interface {
	parseError(error) *parseError
}

type executer interface {
	coordinator
	Exec(*tplExec) *execError
}

type getter interface {
	coordinator
	Get(*tplExec) error
}

type setter interface {
	coordinator
	Set(*tplExec, interface{}) error
}

type position struct {
	tplName   string
	line, pos int
}

func (s *position) parseError(err error) *parseError {
	return &parseError{s.tplName, s.line, s.pos, err}
}

func (s *position) execError(err error) *execError {
	return &execError{s.tplName, s.line, s.pos, err}
}

type execText struct {
	*position
	src []byte
}

func (s execText) Exec(exec *tplExec) *execError {
	return exec.Write(s.src)
}

type execBase struct {
	line, pos int
}

func (s *execBase) Exec() (err error) {
	return fmt.Errorf("ALLLLLLLLLLLLLLLLLLLL")
}
