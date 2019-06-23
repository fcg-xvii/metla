package metla

type cReturn struct {
	position
}

func (s cReturn) execType() execType {
	return execCommand
}

func (s cReturn) exec(exec *tplExec) *execError {
	exec.stack.PopAll()
	exec.returnFlag = true
	return nil
}

type cExit struct {
	position
}

func (s cExit) execType() execType {
	return execCommand
}

func (s cExit) exec(exec *tplExec) *execError {
	exec.stack.PopAll()
	exec.exitFlag = true
	return nil
}
