package metla

import (
	"errors"
	_ "fmt"
	"io"

	"github.com/fcg-xvii/lineman"
)

func init() {
	creators = append(creators, &valueCreator{
		checker:     checkValVariable,
		constructor: newValVariable,
	})
}

func checkValVariable(src []byte) bool {
	if lineman.CheckFirsNameChar(src) > 0 {
		for i := 1; i < len(src); i++ {
			if lineman.CheckBodyNameChar(src[i:]) == 0 {
				return src[i] != '(' && src[i] != '['
			}
		}
	}
	return false
}

// Конструктор строки.
func newValVariable(p *parser) (res token, err error) {
	if name, check := p.ReadName(); !check {
		err = errors.New("Variable parse error :: Unexpected name")
	} else {
		res = &valVariable{string(name)}
	}
	return
}

type valVariable struct {
	name string
}

func (s *valVariable) Val() interface{} {
	return s.name
}

func (s *valVariable) Data(w io.Writer, sto *storage) (err error) {
	return nil
}

func (s *valVariable) String() string     { return "[variable :: {" + s.name + "}]" }
func (s *valVariable) IsExecutable() bool { return false }
