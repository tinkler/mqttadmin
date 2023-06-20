package gsrv

import (
	"net"
	"testing"

	pb_page_v1 "github.com/tinkler/mqttadmin/page/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createServer() (func(), error) {
	s := grpc.NewServer()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return nil, err
	}
	pb_page_v1.RegisterPageGsrvServer(s, NewPageGsrv())
	go s.Serve(lis)
	return func() { s.Stop() }, nil
}

func createClient() (pb_page_v1.PageGsrvClient, func(), error) {
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	c := pb_page_v1.NewPageGsrvClient(conn)
	return c, func() { conn.Close() }, nil
}

func TestGenUser(t *testing.T) {
	closeServer, err := createServer()
	if err != nil {
		t.Fatal(err)
	}
	defer closeServer()
	_, closeClient, err := createClient()
	if err != nil {
		t.Fatal(err)
	}
	defer closeClient()
	// TODO: check method
}
