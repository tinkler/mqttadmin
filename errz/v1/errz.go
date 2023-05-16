package errz_v1

func (m *ValidateError) Error() string {
	return m.Message
}

func (m *AuthError) Error() string {
	return m.Message
}

func (m *ServerError) Error() string {
	return m.Message
}
