package metla

import (
	"fmt"
	_ "io"

	_ "github.com/golang-collections/collections/stack"
)

func newValPrint(p *parser) (err error) {
	p.stack.Push(&execCommand{p.infoRecordFromMark(), execPrint, "print"})
	p.ForwardPos(2)
	for !p.IsEndDocument() {
		p.PassSpaces()
		if p.PosMatchSlice([]byte("}}")) {
			p.ForwardPos(2)
			if p.stack.Len() == 0 {
				err = p.positionError("Empty print tag")
				return
			}
			p.flushStack()
			return
		} else if _, err = initCodeVal(p); err != nil {
			return
		}
	}
	err = p.positionError("Unclosed print tag")
	return
}

func execPrint(exec *tplExec, info *rawInfoRecord) (err error) {
	//fmt.Println("EXEC_PRINT")
	//fmt.Printf("%T, %v\n", exec.st.Peek(), exec.st.Len())
	if exec.st.Len() == 1 {
		_, err = exec.w.Write([]byte(fmt.Sprint(exec.st.Pop())))
	} else if exec.st.Len() > 1 {
		for exec.st.Len() > 0 {
			if _, err = exec.w.Write([]byte(fmt.Sprint(exec.st.Pop()))); err != nil {
				return
			} else if exec.st.Len() > 0 {
				if _, err = exec.w.Write([]byte{',', ' '}); err != nil {
					return
				}
			}
		}
	}
	/*if command, check := exec.st.Peek().(*execCommand); check {
		exec.st.Pop()
		return command.method(exec, command.rawInfoRecord)
	}
	_, err = exec.w.Write([]byte(fmt.Sprint(exec.st.Pop())))*/
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
