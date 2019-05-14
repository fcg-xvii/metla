package prod

import "fmt"

func init() {
	keywords["include"] = newInclude
}

func newInclude(p *parser) *parseError {
	res := include{
		position: position{p.tplName, p.Line(), p.Pos()},
		root:     p.root,
	}
	if err := p.initCodeVal(); err != nil {
		return err
	}
	res.resource = p.stack.Pop().(coordinator)
	p.PassSpaces()
	if !p.IsEndLine() {
		if err := p.initCodeVal(); err != nil {
			return err
		}
		crd := p.stack.Pop().(coordinator)
		if obj, check := crd.(object); !check {
			return crd.parseError("Expected object token")
		} else {
			res.params = obj.vals
		}
	}
	p.stack.Push(res)
	return nil
}

type include struct {
	position
	root     *Metla
	resource coordinator
	params   map[string]coordinator
}

func (s include) execType() execType {
	return execInclude
}

func (s include) exec(exec *tplExec) *execError {
	//fmt.Println("EXEC")
	if iface, err := execOneReturn(s.resource, exec); err != nil {
		return err
	} else if path, check := iface.(string); !check {
		return s.resource.execError("String value expected")
	} else {
		if tpl, check := s.root.template(path); check {
			m := make(map[string]interface{})
			for key, crd := range s.params {
				if iface, err := execOneReturn(crd, exec); err != nil {
					return err
				} else {
					m[key] = iface
				}
			}
			if exists, modified, err := tpl.content(exec.writer, m); err != nil {
				return s.execError(err.Error())
			} else if !exists {
				return s.execError(fmt.Sprintf("Template %v is not exists\n", path))
			} else {
				if exec.modified.Before(modified) {
					exec.modified = modified
				}
			}
		} else {
			return s.execError(fmt.Sprintf("Template %v is not exists\n", path))
		}
	}
	//tpl, check := s.root.template()
	return nil
}
