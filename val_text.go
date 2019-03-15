package metla

import (
	"io"

	"github.com/golang-collections/collections/stack"
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
		p.tpl.pushToken(src)
		p.tpl.pushToken(&execCommand{info, execText})
	}
	return nil
}

func execText(com []interface{}, st *stack.Stack, sto *storage, w io.Writer) (newCom []interface{}, err error) {
	if _, err = w.Write(st.Pop().([]byte)); err == nil {
		newCom = com[1:]
	}
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
