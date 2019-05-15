package metla

import (
	_ "fmt"
)

func init() {
	keywords["if"] = newIf
	keywords["else"] = newElse
}

func findThread(p *parser) (*thread, int, bool) {
	for i := len(p.execList) - 1; i >= 0; i-- {
		if res, check := p.execList[i].(*thread); check {
			if !res.closed {
				return res, i, true
			}
		}
	}
	return nil, 0, false
}

func newIf(p *parser) *parseError {
	//fmt.Println("NEW_IF")
	p.threadLayout++
	p.store.incLayout()
	ck := &thread{position: position{p.tplName, p.Line(), p.LinePos()}}
	for !p.IsEndDocument() {
		if err := p.initCodeVal(); err != nil {
			return err
		}
		p.PassSpaces()
		if p.IsEndLine() {
			crd := p.stack.Pop().(coordinator)
			if pn, check := crd.(rpn); !check {
				return crd.parseError("Expected expression token")
			} else {
				block := &threadBlock{position: ck.position, pn: &pn}
				ck.blocks = append(ck.blocks, block)
				p.stack.Push(ck)
			}
			return nil
		}
	}
	return p.initParseError(p.Line(), p.LinePos(), "check :: expected endline")
}

func newElse(p *parser) *parseError {
	ck, i, check := findThread(p)
	if !check {
		return p.initParseError(p.Line(), p.LinePos(), "Unexpected else token - 'if' token not found")
	} else {
		lastBlock := ck.blocks[len(ck.blocks)-1]
		if lastBlock.pn == nil {
			return lastBlock.parseError("expected 'else if' token before.")
		}
		lastBlock.commands = make([]executer, len(p.execList)-i-1)
		copy(lastBlock.commands, p.execList[i+1:])
		p.execList = p.execList[:i+1]
	}
	block := &threadBlock{position: position{p.tplName, p.Line(), p.LinePos()}}
	p.PassSpaces()
	if p.PosMatchSlice([]byte("if")) {
		p.ForwardPos(2)
		for !p.IsEndLine() {
			if err := p.initCodeVal(); err != nil {
				return err
			}
			p.PassSpaces()
		}
		if p.stack.Peek() == nil {
			return block.parseError("Expected expression token")
		}
		crd := p.stack.Pop().(coordinator)
		if pn, check := crd.(rpn); !check {
			return crd.parseError("Expected expression token")
		} else {
			block.pn = &pn
			ck.blocks = append(ck.blocks, block)
		}
	} else {
		ck.blocks = append(ck.blocks, block)
	}
	return nil
}

type threadBlock struct {
	position
	pn       *rpn
	commands []executer
}

func (s *threadBlock) exec(exec *tplExec) *execError {
	for _, v := range s.commands {
		if err := v.exec(exec); err != nil {
			return err
		}
	}
	return nil
}

type thread struct {
	position
	blocks []*threadBlock
	closed bool
}

func (s *thread) execType() execType {
	return execIf
}

func (s *thread) exec(exec *tplExec) *execError {
	var check bool
	for _, v := range s.blocks {
		if v.pn != nil {
			if err := v.pn.execToBool(exec, &check); err != nil {
				return err
			}
			if check {
				return v.exec(exec)
			}
		} else {
			return v.exec(exec)
		}
	}
	return nil
}
