package metla

import "fmt"

func newValCode(p *parser) error {
	p.ForwardPos(2)
	for !p.IsEndDocument() {
		p.PassSpaces()
		switch p.Char() {
		case ';', '\n':
			p.flushStack()
			p.IncPos()
		case '%':
			if p.NextChar() == '}' {
				p.flushStack()
				p.ForwardPos(2)
				return nil
			}
		default:
			if _, err := initCodeVal(p); err != nil {
				return err
			}
		}
	}
	return p.positionError("CODE_ERR")
}

func newValSet(p *parser) (res interface{}, err error) {
	fmt.Println("INIT_SET")
	info := p.infoRecordFromPos()
	if p.stack.Len() == 0 {
		err = p.positionError("left side of set is empty")
		return
	}
	varsCount, storeUpdate := 0, true
	peekVar := func() (err error) {
		if _, check := p.stack.Peek().(*valVariable); !check {
			err = p.positionError("Variable expected")
		} else {
			varsCount++
		}
		return
	}

	if err = peekVar(); err != nil {
		return
	}

	tmp := p.stack.Pop()
	p.stack.Push(&execMarker{"endset"})
	p.stack.Push(tmp)

	if p.Char() != '=' {
		p.PassSpaces()
		for !p.IsEndDocument() && p.Char() != '=' {
			if p.Char() == ':' {
				if p.NextChar() != '=' {
					err = fmt.Errorf("Expected character '='")
					return
				}
				storeUpdate = false
				p.IncPos()
				break
			}
			if p.Char() != ',' {
				err = p.positionError("Expected '=' or ',' character")
				return
			}
			p.IncPos()
			if _, err = initCodeVal(p); err != nil {
				return
			} else if err = peekVar(); err != nil {
				return
			}
			p.PassSpaces()
		}
	} else {
		storeUpdate = !(p.PrevChar() == ':')
	}
	p.IncPos()
	//p.PassSpaces()
	for !p.IsEndDocument() {
		if _, err = initCodeVal(p); err != nil {
			return
		}
		p.PassSpaces()
		//fmt.Println("SetVal", p.EndLineContent(), p.Char(), p.IsEndLine())
		if p.IsEndLine() {
			p.stack.Push(varsCount)
			p.stack.Push(storeUpdate)
			res = &execCommand{info, execSet, 0}
			p.stack.Push(res)
			return
		}
		if p.Char() != ',' {
			err = p.positionError("Expected ',' or endline character")
			return
		}
		p.IncPos()
	}
	err = p.positionError("Unexpected end of document")
	return
}

func execSet(exec *tplExec, info *rawInfoRecord) (err error) {
	storeUpdate, varsCount := exec.st.Pop().(bool), exec.st.Pop().(int)
	fmt.Println("STORE_UPADTE", storeUpdate)
	var args []interface{}
loop:
	for exec.st.Len() > 0 {
		val := exec.st.Pop()
		switch val.(type) {
		case *execMarker:
			break loop
		default:
			args = append(args, val)
		}
	}
	if len(args) != varsCount*2 {
		err = info.fatalError(fmt.Sprintf("Mismatch count variables and values [%v vs %v]", varsCount, len(args)-varsCount))
	}
	fmt.Println(args)
	l := int(len(args) / 2)
	for i := 0; i < l; i++ {
		v, val := args[i+l].(*variable), args[i]
		v.value = val
		if storeUpdate {
			exec.sto.updateVariable(v)
		} else {
			exec.sto.appendVariable(v)
		}
	}
	return
}

func newValIndex(p *parser) (interface{}, error) {
	return nil, p.positionError("index_error")
}

func newValField(p *parser) (interface{}, error) {
	return nil, p.positionError("field_error")
}
