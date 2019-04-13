package prod

import (
	"io"

	"github.com/fcg-xvii/containers"
)

type tplExec struct {
	tplName  string
	execList []executer
	writer   io.Writer
	sto      *execStorage
	stack    *containers.Stack
}

func (s *tplExec) Write(data []byte) *execError {
	_, err := s.writer.Write(data)
	if err != nil {
		return &execError{s.tplName, 0, 0, err}
	}
	return nil
}

func (s *tplExec) exec() error {
	for _, v := range s.execList {
		if err := v.Exec(s); err != nil {
			return err
		}
	}
	return nil
}
