package logclickhouse

import (
	"database/sql"
	"hezzltask5/internal/models"

	_ "github.com/ClickHouse/clickhouse-go"
)

type LogClickHouse struct {
	DB *sql.DB
}

func NewClickHouseClient(addr string) (*sql.DB, error) {
	db, err := sql.Open("clickhouse", addr)
	if err != nil {

		return nil, err
	}
	return db, nil

}

func NewLogClickHouse(db *sql.DB) *LogClickHouse {
	return &LogClickHouse{DB: db}
}

func (lc *LogClickHouse) Upload(batch []models.Event) error {
	data := make([]models.Event, len(batch))
	copy(data, batch)

	tx, err := lc.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()

		}
	}()

	stmt, err := tx.Prepare("INSERT INTO events (id, campaign_id, name,description,priority,removed,created_at) VALUES (?, ?, ?,?,?,?,?)")
	if err != nil {
		tx.Rollback()
		return err

	}
	defer stmt.Close()
	var boolValue int
	for _, elem := range data {

		boolValue = 0
		if elem.Removed {
			boolValue = 1
		}
		_, err := stmt.Exec(elem.Id, elem.CampaignId, elem.Name, elem.Description, elem.Priority, boolValue, elem.EventTime)
		if err != nil {

			return err
		}

	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil

}
