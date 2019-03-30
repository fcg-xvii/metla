package metla

import (
	"fmt"
)

func init() {
	keywords["if"] = keywordIf
	keywords["else"] = keywordElse
	keywords["endif"] = keywordEndif
}

func keywordIf(p *parser) (res interface{}, err error) {
	res = &execCommand{p.infoRecordFromPos(), execIf, 0}
	p.stack.Push(&execCommand{p.infoRecordFromPos(), execIf, 0})
	p.openStack.Push(openFlag{p.infoRecordFromPos(), "if"})
	for !p.IsEndLine() {
		if _, err = initCodeVal(p); err != nil {
			return
		}
	}
	return
}

func keywordElse(p *parser) (res interface{}, err error) {
	if p.openStack.Len() == 0 {
		err = p.positionError("else without opened condition")
		return
	} else {
		if openInfo := p.openStack.Peek().(openFlag); openInfo.tagName != "if" {
			err = openInfo.info.fatalError(fmt.Sprintf("condition close with unclosed %v tag", openInfo.tagName))
			return
		}
	}
	p.stack.Push(&execMarker{"else"})
	return
}

func keywordEndif(p *parser) (res interface{}, err error) {
	if p.openStack.Len() == 0 {
		err = p.positionError("endif without opened condition")
		return
	} else {
		if openInfo := p.openStack.Pop().(openFlag); openInfo.tagName != "if" {
			err = openInfo.info.fatalError(fmt.Sprintf("condition close with unclosed %v tag", openInfo.tagName))
			return
		}
	}
	p.stack.Push(&execMarker{"endif"})
	return
}

func convertCondition(exec *tplExec, info *rawInfoRecord) (err error) {
	condIface := exec.st.Pop()
	if vlr, check := condIface.(valuer); check {
		condIface = vlr.Value()
	}
	if _, check := convert(condIface, true); !check {
		err = info.fatalError("Expected boolean value")
		return
	}
	exec.st.Push(condIface)
	return
}

func execCondition(exec *tplExec) (err error) {
	fmt.Println("EXEC_CONDITION....................", exec.st.Peek().(bool), exec.index)
	condition := exec.st.Pop().(bool)
	if condition {
		for {
			if err = exec.execNext(); err != nil {
				return
			} else if m, check := exec.st.Peek().(*execMarker); check {
				fmt.Println("NAME.....", m.name)
				switch m.name {
				case "else":
					exec.st.Pop()
					exec.st.Push(!condition)
					execCondition(exec)
					return
				case "endif":
					exec.st.Pop()
					return

				}
			}
		}
	} else {
		for _, v := range exec.list[exec.index:] {
			if marker, check := v.(*execMarker); check {
				switch marker.name {
				case "else":
					//exec.st.Push(!condition)
					execCondition(exec)
					return
				case "endif":
					return
				}
				if marker.name == "endif" {
					return
				}
			}
			exec.index++
		}
	}
	return
}

func execIf(exec *tplExec, info *rawInfoRecord) (err error) {
	if err = convertCondition(exec, info); err == nil {
		err = execCondition(exec)
	}
	return
}
