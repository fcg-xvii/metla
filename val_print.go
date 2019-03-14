package metla

import (
	"fmt"
	"io"
)

func newValPrint(p *parser, parent tokenContainer) (err error) {
	p.SetupMark()
	p.ForwardPos(2)
	for !p.IsEndDocument() {
		p.PassSpaces()
		if err = initCodeVal(p, nil); err != nil {
			return
		} else {
			p.PassSpaces()
			if p.PosMatchSlice([]byte("}}")) {
				fmt.Println("PEEK......", p.codeStack.Peek())
				if p.codeStack.Peek() == nil {
					err = p.positionError("Empty print tag")
				} else {
					p.flushToParent(parent, &valPrint{p.infoRecordFromMark(), p.codeStack.Pop().(token)})
					p.ForwardPos(2)
				}
			} else if !p.IsEndDocument() {
				p.SetupMark()
				err = p.positionError(fmt.Sprintf("Unexpected char '%c'", p.Char()))
			} else {
				err = p.positionError("Unclosed print tag")
			}
			return
		}
	}
	return
}

type valPrint struct {
	*rawInfoRecord
	t token
}

func (s *valPrint) IsExecutable() bool { return true }
func (s *valPrint) String() string     { return "[ print { " + s.t.String() + "} ]" }

func (s *valPrint) execObject(sto *storage, tpl *template, parent executor) (res executor, err error) {
	if res, err = s.t.execObject(sto, tpl, parent); err == nil {
		res = &printExec{s.rawInfoRecord, res}
	}
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
}
