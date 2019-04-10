package metla

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
		p.stack.Push(&execCommand{info, execText, "text", nil})
		p.stack.Push(src)
		p.flushStack()
	}
	return nil
}

func execText(exec *tplExec, info *rawInfoRecord) (err error) {
	_, err = exec.w.Write(exec.st.Pop().([]byte))
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
