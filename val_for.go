package metla

import "fmt"

func init() {
	keywords["for"] = keywordFor
	keywords["endfor"] = keywordEndfor
}

func keywordFor(p *parser) (res interface{}, err error) {
	p.PassSpaces()
	p.stack.Push(&execCommand{p.infoRecordFromPos(), execFor, 0})
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
	p.stack.Push(&execMarker{"endfor"})
	return
}

func execFor(exec *tplExec, info *rawInfoRecord) (err error) {
	indexVar, minVal, maxVal := exec.st.Pop().(*variable), exec.st.Pop().(int64), exec.st.Pop().(int64)
	codePos := exec.index
	indexVar.value = minVal
	exec.sto.newLayout()
	exec.sto.appendVariable(indexVar)
	for indexVar.value.(int64) < maxVal {
		//fmt.Println("!!!")
		exec.index = codePos
		for {
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
