package metla

import (
	"fmt"
)

func newValText(p *parser) error {
	p.SetupMark()
	info := p.infoRecordFromMark()
	for !p.IsEndDocument() {
		if cur, next := p.Char(), p.NextChar(); cur == '{' && (next == '{' || next == '*' || next == '%') {
			break
		} else {
			p.IncPos()
		}
	}
	if src := p.MarkVal(0); len(src) > 0 {
		//p.tpl.pushToken(src)
		//p.tpl.pushToken(&execCommand{info, execText})
		p.stack.Push(src)
		p.stack.Push(&execCommand{info, execText, 2})
	}
	return nil
}

func execText(exec *tplExec, info *rawInfoRecord) (err error) {
	fmt.Println("EXEC_TEXT", exec.st.Len())
	_, err = exec.w.Write(exec.st.Pop().([]byte))
	fmt.Println("ERRRRR", err, exec.st.Len())
	return
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func newValComment(p *parser) error {
	p.SetupMark()
	p.ForwardPos(2)
	for !p.IsEndDocument() {
		if p.Char() == '*' && p.NextChar() == '}' {
			p.ForwardPos(2)
			return nil
		}
		p.IncPos()
	}
	return p.positionError("Unclosed comment tag")
}
