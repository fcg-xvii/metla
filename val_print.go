package metla

import (
	"fmt"
	"io"

	"github.com/golang-collections/collections/stack"
)

func newValPrint(p *parser) (err error) {
	p.SetupMark()
	info := p.infoRecordFromMark()
	stackOffset := p.stack.Len()
	p.ForwardPos(2)
	for !p.IsEndDocument() {
		p.PassSpaces()
		if p.PosMatchSlice([]byte("}}")) {
			p.ForwardPos(2)
			if p.stack.Len() == 0 {
				err = p.positionError("Empty print tag")
				return
			} else if exec, check := p.stack.Peek().(*execCommand); check {
				if exec.itemsCount < p.stack.Len()-stackOffset {
					err = p.stack.Peek().(token).fatalError("More one value")
					return
				}
			} else if p.stack.Len()-stackOffset > 1 {
				err = info.fatalError("More one value")
				return
			}
			p.stack.Push(&execCommand{info, execPrint, p.stack.Len() + 1})
			return
		} else if _, err = initCodeVal(p); err != nil {
			return
		}
	}
	err = p.positionError("Unclosed print tag")
	return
}

func execPrint(com []interface{}, st *stack.Stack, sto *storage, w io.Writer) (newCom []interface{}, err error) {
	if _, err = w.Write([]byte(fmt.Sprint(st.Pop()))); err == nil {
		newCom = com[1:]
	}
	return
}

// ======================================================================

/*type valPrint struct {
	*rawInfoRecord
}

func (s *valPrint) IsExecutable() bool { return true }
func (s *valPrint) String() string     { return "[ print... ]" }

func (s *valPrint) execObject(sto *storage) (res executor, err error) {

	return
}

type printExec struct {
	*rawInfoRecord
	e executor
}

func (s *printExec) String() string { return "[ print { " + s.e.String() + "} ]" }
func (s *printExec) Data(w io.Writer) (err error) {
	_, err = w.Write([]byte(s.e.String()))
	return
}*/
