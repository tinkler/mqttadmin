package status

import "google.golang.org/grpc/status"

func (s *Status) Is(err error) bool {
	if s == nil {
		return err == nil
	}
	if st, ok := err.(*Status); ok {
		return s.Code == st.Code && s.Message == st.Message
	}
	if ge, ok := status.FromError(err); ok {
		s.ConvToType(1)
		return uint32(s.Code) == uint32(ge.Code())
	}
	return false
}
