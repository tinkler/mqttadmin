package page

import (
	"io"

	"github.com/tinkler/mqttadmin/pkg/gs"
	"github.com/tinkler/mqttadmin/pkg/logger"
)

type PageRow struct {
	RowNo int
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
		err = stream.Send(nil)
		if err != nil {
			logger.Error(err)
			return err
		}
	}
}
