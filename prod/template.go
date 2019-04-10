package prod

import (
	"io"
)

type tplExec struct {
	execList []executer
	writer   io.Writer
	sto      *storage
}

func (s *tplExec) Write(data []byte) error {
	_, err := s.writer.Write(data)
	return err
}

func (s *tplExec) exec() error {
	for _, v := range s.execList {
		if err := v.Exec(s); err != nil {
			return err
		}
	}
	return nil
}
