package core

type core struct {
	errs []error
	err error
}

func (s *core) setErrs(err error)  {
	s.errs = append(s.errs, err)
}

func (s *core) IsErr(err ...error) bool {
	if len(err) == 0 {
		for _,v := range s.errs {
			if v != nil {
				s.err = v
				return true
			}
		}
	} else {
		for _,v := range err {
			if v != nil {
				s.err = v
				return true
			}
		}
	}

	return false
}

func (s *core) GetErr() error {
	return s.err
}