package gsrv

// type PagePageRowGenUserStream struct {
// 	stream pb_page_v1.PageGsrv_PageRowGenRowServer
// 	m      *page.PageRow
// }

// func (s *PagePageRowGenUserStream) Context() context.Context {
// 	return s.stream.Context()
// }
// func (s *PagePageRowGenUserStream) Send(v *user.User) error {
// 	res := mrz.NewTypedRes[*pb_page_v1.PageRow, *pb_user_v1.User]()
// 	// data
// 	res.Data = new(pb_page_v1.PageRow)
// 	jd, err := sjson.Marshal(s.m)
// 	if err != nil {
// 		return err
// 	}
// 	err = sjson.Unmarshal(jd, res.Data)
// 	if err != nil {
// 		return err
// 	}
// 	// resp
// 	respByt, _ := sjson.Marshal(v)
// 	newResp := new(pb_user_v1.User)
// 	err = sjson.Unmarshal(respByt, newResp)
// 	if err != nil {
// 		return err
// 	}
// 	res.Resp = newResp
// 	return s.stream.Send(res.ToAny())
// }
// func (s *PagePageRowGenUserStream) Recv() (*user.User, error) {
// 	in, err := s.stream.Recv()
// 	if err != nil {
// 		return nil, err
// 	}
// 	req := mrz.ToTypedModel[*pb_page_v1.PageRow, *pb_user_v1.User](in)
// 	jd, err := sjson.Marshal(req.Data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = sjson.Unmarshal(jd, s.m)
// 	if err != nil {
// 		return nil, err
// 	}
// 	argsByt, _ := sjson.Marshal(req.Args)
// 	newArgs := new(user.User)
// 	err = sjson.Unmarshal(argsByt, newArgs)
// 	return newArgs, err
// }

// func (s *pageGsrv) PageRowGenUser(stream pb_page_v1.PageGsrv_PageRowGenUserServer) error {
// 	gsStream := &PagePageRowGenUserStream{stream, &page.PageRow{}}
// 	return gsStream.m.GenUser(gsStream)
// }
