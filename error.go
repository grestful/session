package session


type SError map[string]error

func (se SError) Error(sid string) error {
	if e,ok := se[sid]; ok {
		return e
	}

	return nil
}

func (se SError) SetErr(sid string, err error)  {
	se[sid] = err
}

func (se SError) Remove(sid string)  {
	delete(se, sid)
}

