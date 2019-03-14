package metla

import "io"

func newValText(p *parser, parent tokenContainer) error {
	p.SetupMark()
	for !p.IsEndDocument() {
		if cur, next := p.Char(), p.NextChar(); cur == '{' && (next == '{' || next == '*' || next == '%') {
			break
		} else {
			p.IncPos()
		}
	}
	if src := p.MarkVal(0); len(src) > 0 {
		parent.pushToken(&valText{p.infoRecordFromMark(), src})
	}
	return nil
}

type valText struct {
	*rawInfoRecord
	src []byte
}

func (s *valText) IsExecutable() bool { return true }
func (s *valText) String() string     { return "[ text ]" }
func (s *valText) execObject(sto *storage, tpl *template, parent executor) (res executor, err error) {
	return s, nil
}

func (s *valText) Data(w io.Writer) (err error) {
	_, err = w.Write(s.src)
	return
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func newValComment(p *parser, parent tokenContainer) error {
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
