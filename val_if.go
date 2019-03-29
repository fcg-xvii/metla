package metla

import (
	"fmt"
)

func init() {
	keywords["if"] = keywordIf
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

func execIf(exec *tplExec, info *rawInfoRecord) (err error) {
	condIface := exec.st.Pop()
	if vlr, check := condIface.(valuer); check {
		condIface = vlr.Value()
	}
	if _, check := convert(condIface, true); !check {
		err = info.fatalError("Expected boolean value")
		return
	}
	if condIface.(bool) {
		for {
			//fmt.Println(exec.index, len(exec.list))

			if err = exec.execNext(); err != nil {
				return
			} else if m, check := exec.st.Peek().(*execMarker); check && m.name == "endif" {
				exec.st.Pop()
				return
			}
		}
	} else {
		for _, v := range exec.list[exec.index:] {
			if marker, check := v.(*execMarker); check {
				if marker.name == "endif" {
					return
				}
			}
			exec.index++
		}
	}
	err = info.fatalError("IF_ERR")
	return
}
