package gsrv

import (
	"context"
	"net"
	"testing"

	mrz_v1 "github.com/tinkler/mqttadmin/mrz/v1"
	pb_user_v1 "github.com/tinkler/mqttadmin/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestUserSrvUserAuthSignin(t *testing.T) {
	s := grpc.NewServer()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	pb_user_v1.RegisterUserGsrvServer(s, NewUserGsrv())
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
	c := pb_user_v1.NewUserGsrvClient(conn)
	req := mrz_v1.NewTypedModel[*pb_user_v1.Auth, *pb_user_v1.Auth]()
	_, err = c.AuthSignin(context.Background(), req.ToAny())
	if err != nil {
		t.Fatal(err)
	}
}
