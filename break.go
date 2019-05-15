package metla

func init() {
	keywords["break"] = newBreak
}

func newBreak(p *parser) *parseError {
	p.stack.Push(keyBreak{position{p.tplName, p.Line(), p.LinePos()}})
	return nil
}

type keyBreak struct {
	position
}

func (s keyBreak) exec(exec *tplExec) *execError {
	exec.breakFlag = true
	return nil
}

func (s keyBreak) execType() execType {
	return execBreak
}
