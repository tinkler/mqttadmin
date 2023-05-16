package status

func (s *Status) Is(err error) bool {
	if s == nil {
		return err == nil
	}
	st, ok := err.(*Status)
	if !ok {
		return false
	}
	return s.Code == st.Code && s.Message == st.Message
}
