// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT.
package gsrv
import (
	"context"
	mrz "github.com/tinkler/mqttadmin/mrz/v1"
	"github.com/tinkler/mqttadmin/pkg/model/page"
	pb_page_v1 "github.com/tinkler/mqttadmin/page/v1"
	anypb "google.golang.org/protobuf/types/known/anypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	"github.com/tinkler/mqttadmin/pkg/gs"
	"github.com/tinkler/mqttadmin/pkg/jsonz/sjson"
	"github.com/tinkler/mqttadmin/pkg/model/user"
)


type pageGsrv struct {
	pb_page_v1.UnimplementedPageGsrvServer
}

func NewPageGsrv() *pageGsrv {
	return &pageGsrv{}
}


func (u *pageGsrv) PageFetchUser(ctx context.Context, in *anypb.Any) (out *anypb.Any, err error) {
	gm := mrz.ToTypedModel[*pb_page_v1.Page, *structpb.Struct](in)
	m := mrz.GetData[*page.Page](gm)
	res := mrz.NewTypedRes[*pb_page_v1.Page, *structpb.Value]()
	var resData []*user.User
	resData, err = m.FetchUser(ctx)
	if err != nil {
		return nil, err
	}
	mrz.SetResData(res, m)
	mrz.SetResRespList(res, resData)
	return res.ToAny(), nil
}
type PagePageRowGenRowStream struct {
	stream pb_page_v1.PageGsrv_PageRowGenRowServer
	m      *page.PageRow
}
func (s *PagePageRowGenRowStream) Context() context.Context {
	return s.stream.Context()
}
func (s *PagePageRowGenRowStream) Send(_ *gs.Null) error {
	res := mrz.NewTypedRes[*pb_page_v1.PageRow, *anypb.Any]()
	// data
	res.Data = new(pb_page_v1.PageRow)
	jd, err := sjson.Marshal(s.m)
	if err != nil {
		return err
	}
	err = sjson.Unmarshal(jd, res.Data)
	if err != nil {
		return err
	}
	// resp
	null := structpb.NewNullValue()
	resp, _ := anypb.New(null)
	res.Resp = resp
	return s.stream.Send(res.ToAny())
}
func (s *PagePageRowGenRowStream) Recv() (*gs.Null, error) {
	in, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}
	req := mrz.ToTypedModel[*pb_page_v1.PageRow, *anypb.Any](in)
	jd, err := sjson.Marshal(req.Data)
	if err != nil {
		return nil, err
	}
	err = sjson.Unmarshal(jd, s.m)
	return nil, err
}

func (s *pageGsrv) PageRowGenRow(stream pb_page_v1.PageGsrv_PageRowGenRowServer) error {
	gsStream := &PagePageRowGenRowStream{stream, &page.PageRow{} }
	return gsStream.m.GenRow(gsStream)
}
