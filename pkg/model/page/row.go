package page

import (
	"io"

	"github.com/tinkler/mqttadmin/pkg/gs"
	"github.com/tinkler/mqttadmin/pkg/logger"
	"github.com/tinkler/mqttadmin/pkg/model/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PageRow struct {
	RowNo    int
	Chapters map[string]Chapter
	Option   map[string]interface{}
}

type Chapter struct {
	Index int
	Name  string
}

// @stream(bidi)
func (m *PageRow) GenRow(stream gs.NullStream) error {
	for {
		_, err := stream.Recv()
		logger.Info("PageRow.GenRow", "err", err)
		if err == io.EOF {
			return stream.Send(nil)
		}
		if err != nil {
			logger.Error(err)
			return err
		}
		m.RowNo++
		if m.RowNo == 10 {
			return status.New(codes.Canceled, "reach 10").Err()
		}
		err = stream.Send(nil)
		if err != nil {
			logger.Error(err)
			return err
		}
	}
}

func (m *PageRow) GenUser(stream gs.Stream[*user.User]) error {
	return nil
}
