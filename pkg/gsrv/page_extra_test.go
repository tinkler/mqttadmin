package gsrv

import (
	"context"
	"net"
	"testing"
	"time"

	mrz "github.com/tinkler/mqttadmin/mrz/v1"
	pb_page_v1 "github.com/tinkler/mqttadmin/page/v1"
	"github.com/tinkler/mqttadmin/pkg/logger"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestPageSrvPageRowGenRow(t *testing.T) {
	s := grpc.NewServer()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	pb_page_v1.RegisterPageGsrvServer(s, NewPageGsrv())
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Error(err)
		}
	}()
	defer s.Stop()

	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := pb_page_v1.NewPageGsrvClient(conn)
	stream, err := c.PageRowGenRow(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	data := &pb_page_v1.PageRow{}
	go func() {
		for {
			any, err := stream.Recv()
			if err != nil && status.Code(err) == codes.Code(code.Code_CANCELLED) {
				break
			}
			if err != nil {
				logger.Error(err)
				continue
			}

			res := mrz.ToDataTypedRes[*pb_page_v1.PageRow](any)
			t.Logf("%v", res.Data.RowNo)
			data.RowNo = res.Data.RowNo
		}
	}()
	for range time.NewTicker(time.Millisecond * 100).C {
		req := mrz.NewDataTypedModel[*pb_page_v1.PageRow]()
		req.Data = data
		if data.RowNo == 10 {
			err := stream.CloseSend()
			if err != nil {
				logger.Error(err)
			}
			break
		}
		err := stream.Send(req.ToAny())
		if err != nil {
			logger.Error(err)
		}

	}
}
