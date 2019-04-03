package metla

import "fmt"

func init() {
	keywords["for"] = keywordFor
	keywords["endfor"] = keywordEndfor
}

func keywordFor(p *parser) (res interface{}, err error) {
	p.PassSpaces()
	p.stack.Push(&execCommand{p.infoRecordFromPos(), execFor, "for"})
	p.openStack.Push(openFlag{p.infoRecordFromPos(), "for"})
	// Получаем переменную индекса
	if res, err = initCodeVal(p); err == nil {
		if _, check := res.(*valVariable); !check {
			err = p.positionError("For index must be a variable")
			return
		}
		p.PassSpaces()
		if !p.PosMatchSlice([]byte("in")) {
			err = p.positionError("Expected in article after for index (for var in min:max[:step]) template")
			return
		}
		p.ForwardPos(2)
		p.PassSpaces()
		if res, err = initCodeVal(p); err == nil {
			p.PassSpaces()
			if p.Char() != ':' {
				err = p.positionError("Expected : after min value (for var in min:max[:step]) template")
				return
			}
			p.IncPos()
			p.PassSpaces()
			if res, err = initCodeVal(p); err == nil {
				p.PassSpaces()
				if p.Char() != ';' && p.Char() != '\n' && !p.PosMatchSlice([]byte("%}")) {
					err = fmt.Errorf("Endline expected")
				}
			}
		}
	}
	//err = p.positionError("ENDFOR")
	return
}

func keywordEndfor(p *parser) (res interface{}, err error) {
	if p.openStack.Len() == 0 {
		err = p.positionError("endfor without opened cycle")
		return
	} else {
		if openInfo := p.openStack.Pop().(openFlag); openInfo.tagName != "for" {
			err = openInfo.info.fatalError(fmt.Sprintf("for close with unclosed %v tag", openInfo.tagName))
			return
		}
	}
	p.stack.Push(&execMarker{"endfor"})
	return
}

func execFor(exec *tplExec, info *rawInfoRecord) (err error) {
	indexVar := exec.st.Pop().(*variable)
	iMinVal, check := convert(exec.st.Pop(), int64(0))
	if !check {
		return fmt.Errorf("Loop min value must be integer")
	}
	iMaxVal, check := convert(exec.st.Pop(), int64(0))
	if !check {
		return fmt.Errorf("Loop max value must be integer")
	}
	minVal, maxVal := iMinVal.(int64), iMaxVal.(int64)
	codePos := exec.index
	indexVar.value = minVal
	exec.sto.newLayout()
	exec.sto.appendVariable(indexVar)
	for indexVar.value.(int64) < maxVal {
		exec.index = codePos
		for {
			//fmt.Println(exec.index, len(exec.list))
			if err = exec.execNext(); err != nil {
				return
			} else if _, check := exec.st.Peek().(*execMarker); check {
				exec.st.Pop()
				break
			}
		}
		indexVar.value = indexVar.value.(int64) + 1
	}
	exec.sto.dropLayout()
	//fmt.Println("OKO", indexVar, minVal, maxVal)*/
	//err = info.fatalError("AAAAAAAAAAAAAAAAAAAAAAA")
	return
}