package metla

import "time"
import _ "fmt"

func init() {
	keywords["include"] = keywordInclude
}

func keywordInclude(p *parser) (res interface{}, err error) {
	res = &execCommand{p.infoRecordFromPos(), execInclude, "include", nil}
	p.stack.Push(res)
	if _, err = initCodeVal(p); err != nil {
		return
	}
	p.PassSpaces()
	if p.Char() != '\n' && p.Char() != '*' {
		var obj interface{}
		if obj, err = initCodeVal(p); err != nil {
			return
		}
		if ex, mCheck := obj.(*execCommand); mCheck && ex.name != "init-object" {
			err = p.positionError("Expected object variable")
			return
		}
	} else {
		p.stack.Push(nil)
	}
	return
}

func execInclude(exec *tplExec, info *rawInfoRecord) (err error) {
	//fmt.Println("EXEC_INCLUDE")
	nameObj, paramsObj, tplName := exec.st.Pop(), exec.st.Pop(), ""
	if nameVar, check := nameObj.(*variable); check {
		nameObj = nameVar.value
	}
	switch nameObj.(type) {
	case string:
		tplName = nameObj.(string)
	default:
		err = info.fatalError("template name string or variable expected")
		return
	}
	if paramsObj != nil {
		layout := exec.sto.newLayout()
		m := paramsObj.(map[string]interface{})
		for key, val := range m {
			if v, check := val.(*variable); check {
				val = v.value
			}
			layout.appendVariable(&variable{key, val, true})
		}
	}
	var tpl *Template
	if tpl, err = exec.root.Template(tplName); err == nil {
		if err = tpl.checkUpdate(); err == nil {
			var modified time.Time
			if modified, err = tpl.result(exec.sto, exec.w); err == nil && modified.After(exec.modified) {
				//fmt.Println("EM", exec.modified, "MM", modified, modified.After(exec.modified))
				exec.modified = modified
			}
		}
	}
	if paramsObj != nil {
		exec.sto.dropLayout()
	}
	return
}
