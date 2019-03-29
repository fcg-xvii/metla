package metla

import (
	"fmt"
)

func newValIndex(p *parser) (res interface{}, err error) {
	p.IncPos()
	tmp := p.stack.Pop()
	res = &execCommand{p.infoRecordFromMark(), execIndex, 0}
	p.stack.Push(res)
	p.stack.Push(tmp)
	for !p.IsEndDocument() {
		p.PassSpaces()
		if p.Char() == ']' {
			p.IncPos()
			//p.stack.Push(res)
			return
		} else {
			if _, err = initCodeVal(p); err != nil {
				return
			}
		}
	}
	err = fmt.Errorf("Expected close index char ']'")
	return
}

func execIndex(exec *tplExec, info *rawInfoRecord) (err error) {
	if _, check := exec.st.Peek().(*variable); !check {
		return info.fatalError("Expected variable")
	}
	obj, index := exec.st.Pop(), exec.st.Pop()
	if v, check := index.(*variable); check {
		index = v.value
	}
	exec.st.Push(indexVariable{obj.(*variable), index})

	return
}
