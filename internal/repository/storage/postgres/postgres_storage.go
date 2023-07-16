package postgres

import (
	"context"
	"hezzltask5/internal/models"
	projerrors "hezzltask5/internal/models/errors"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	createItem = `
	insert into items(campaign_id, name)
	 values ($1, $2)
	 returning id,campaign_id,name,description,priority,removed,created_at`
	updateItem1Part = `
	select id, campaign_id, name, description, priority, removed, created_at 
	from items
	where id = $1 and campaign_id = $2
	for update
`
	updateItem2Part = `
		update items 
        set name = $1, description = $2
        where id = $3 and campaign_id = $4
		returning id,campaign_id,name,description,priority,removed,created_at
`
	markItemAsDelited = `
		update items 
        set removed = true
        where id = $1 and campaign_id = $2
		returning id,campaign_id,name,description,priority,removed,created_at
`
	getList = `
select * from items where removed<>true 
`
)

type PostgresStorage struct {
	DB *pgx.Conn
	sync.Mutex
}

func NewPostgresStorage(conn *pgx.Conn) *PostgresStorage {

	return &PostgresStorage{
		DB: conn,
	}
}

func NewPostgresClient(dbsn string) (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	conn, err := pgx.Connect(ctx, dbsn)

	if err != nil {

		return nil, err
	}
	return conn, err
}

func (ps *PostgresStorage) Screate(ctx context.Context, campaignId int, name string) (models.Item, error) {
	ps.Lock()
	defer ps.Unlock()
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	var item models.Item

	if err := ps.DB.QueryRow(ctx, createItem, campaignId, name).Scan(&item.Id, &item.CampaignId, &item.Name, &item.Description, &item.Priority, &item.Removed, &item.CreatedAt); err != nil {

		return models.Item{}, nil
	}

	return item, nil
}

func (ps *PostgresStorage) Supdate(ctx context.Context, id, campaignId int, name, description string) (models.Item, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	tx, err := ps.DB.Begin(ctx)
	if err != nil {
		return models.Item{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())

		} else {
			tx.Commit(context.TODO())
		}
	}()
	var item models.Item
	if err := tx.QueryRow(ctx, updateItem1Part, id, campaignId).Scan(&item.Id, &item.CampaignId, &item.Name, &item.Description, &item.Priority, &item.Removed, &item.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return models.Item{}, projerrors.ErrItemNotFound
		}
		return models.Item{}, err
	}

	if item.Removed {
		return models.Item{}, projerrors.ErrItemNotFound
	}

	if err := ps.DB.QueryRow(ctx, updateItem2Part, name, description, id, campaignId).Scan(&item.Id, &item.CampaignId, &item.Name, &item.Description, &item.Priority, &item.Removed, &item.CreatedAt); err != nil {

		return models.Item{}, err
	}

	return item, nil
}

func (ps *PostgresStorage) Sdelete(ctx context.Context, id, campaignId int) (models.Item, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	tx, err := ps.DB.Begin(ctx)
	if err != nil {
		return models.Item{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())

		} else {
			tx.Commit(context.TODO())
		}
	}()
	var record models.Item
	if err := tx.QueryRow(ctx, updateItem1Part, id, campaignId).Scan(&record.Id, &record.CampaignId, &record.Name, &record.Description, &record.Priority, &record.Removed, &record.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return models.Item{}, projerrors.ErrItemNotFound
		}
		return models.Item{}, err
	}

	if record.Removed {
		return models.Item{}, projerrors.ErrItemNotFound
	}

	var item models.Item
	if err := ps.DB.QueryRow(ctx, markItemAsDelited, id, campaignId).Scan(&item.Id, &item.CampaignId, &item.Name, &item.Description, &item.Priority, &item.Removed, &item.CreatedAt); err != nil {

		return models.Item{}, err
	}

	return item, nil

}

func (ps *PostgresStorage) SgetList(ctx context.Context) ([]models.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	list := make([]models.Item, 0)

	rows, err := ps.DB.Query(ctx, `select * from items where removed<>true ;`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var item models.Item
		err = rows.Scan(&item.Id, &item.CampaignId, &item.Name, &item.Description, &item.Priority, &item.Removed, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return list, nil
}
