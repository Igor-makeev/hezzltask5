package logclickhouse

import (
	"context"
	"fmt"
	"hezzltask5/internal/models"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

const insertQuery = "INSERT INTO events"

type LogClickHouse struct {
	conn driver.Conn
}

func NewClickHouseClient(addr string) (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{addr},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil

}

func NewLogClickHouse(conn driver.Conn) *LogClickHouse {
	return &LogClickHouse{conn: conn}
}

func (lc *LogClickHouse) Upload(data []models.Event) error {

	localData := make([]models.Event, len(data))
	copy(localData, data)

	batch, err := lc.conn.PrepareBatch(context.Background(), insertQuery)
	if err != nil {
		return err
	}

	for _, elem := range data {

		err := batch.Append(
			int32(elem.Id),
			int32(elem.CampaignId),
			elem.Name,
			elem.Description,
			int32(elem.Priority),
			elem.Removed,
			elem.EventTime,
		)
		if err != nil {
			return err
		}
	}
	return batch.Send()

}
