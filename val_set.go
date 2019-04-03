package metla

import "fmt"

func newValSet(p *parser) (res interface{}, err error) {
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
	res = &execCommand{info, execSet, "set"}
	p.stack.Push(res)
	p.stack.Push(tmp)
	p.pushSplitter()

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
			p.pushSplitter()
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
	p.stack.Push(&execMarker{"endvars"})
	p.stack.Push(storeUpdate)
	//p.pushSplitter()
	for !p.IsEndDocument() {
		p.PassSpaces()
		if p.IsEndLine() {
			p.stack.Push(&execMarker{"endset"})
			return
		} else if p.Char() == ',' {
			p.pushSplitter()
			p.IncPos()
		} else {
			if _, err = initCodeVal(p); err != nil {
				return
			}
		}
	}
	err = p.positionError("Unexpected end of document")
	return
}

func execSet(exec *tplExec, info *rawInfoRecord) (err error) {
	endVarsAccepted, storeUpdate, varsCount := false, false, 0
	//fmt.Println("STORE_UPADTE", storeUpdate)
	var args []interface{}
loop:
	for exec.st.Len() > 0 {
		val := exec.st.Pop()
		switch val.(type) {
		case *execMarker:
			if endVarsAccepted {
				break loop
			} else {
				storeUpdate, varsCount, endVarsAccepted = exec.st.Pop().(bool), len(args), true
			}
		default:
			args = append(args, val)
		}
	}
	if len(args) != varsCount*2 {
		err = info.fatalError(fmt.Sprintf("Mismatch count variables and values [%v vs %v]", varsCount, len(args)-varsCount))
	}
	//fmt.Println(args)
	l := int(len(args) / 2)
	for i := 0; i < l; i++ {
		v, val := args[i].(*variable), args[i+l]
		v.value = val
		if storeUpdate {
			//fmt.Println("store-update")
			exec.sto.updateVariable(v)
		} else {
			//fmt.Println("store-append")
			exec.sto.appendVariable(v)
		}
	}
	return
}
