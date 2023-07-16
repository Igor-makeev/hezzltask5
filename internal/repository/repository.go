package repository

import (
	"context"
	"hezzltask5/internal/models"

	"github.com/sirupsen/logrus"
)

type persistanceStorage interface {
	Screate(ctx context.Context, campaignId int, name string) (models.Item, error)
	Supdate(ctx context.Context, id, campaignId int, name, description string) (models.Item, error)
	Sdelete(ctx context.Context, id, campaignId int) (models.Item, error)
	SgetList(ctx context.Context) ([]models.Item, error)
}

type cache interface {
	Get(ctx context.Context) ([]models.Item, error)
	Set(ctx context.Context, items []models.Item)
	Remove(ctx context.Context)
}

type Repository struct {
	persistanceStorage
	cache
}

func NewRepository(ps persistanceStorage, c cache) *Repository {
	return &Repository{persistanceStorage: ps, cache: c}
}

func (r *Repository) Create(ctx context.Context, campaignId int, name string) (models.Item, error) {
	return r.persistanceStorage.Screate(ctx, campaignId, name)
}

func (r *Repository) Update(ctx context.Context, id, campaignId int, name, description string) (models.Item, error) {
	r.cache.Remove(ctx)
	return r.persistanceStorage.Supdate(ctx, id, campaignId, name, description)
}

func (r *Repository) Delete(ctx context.Context, id, campaignId int) (models.Item, error) {
	r.cache.Remove(ctx)
	return r.persistanceStorage.Sdelete(ctx, id, campaignId)
}

func (r *Repository) GetList(ctx context.Context) ([]models.Item, error) {

	if data, err := r.cache.Get(ctx); data != nil {
		return data, nil
	} else {
		logrus.Printf("repo level: getlist method: %v", err)
	}
	data, err := r.persistanceStorage.SgetList(ctx)
	if err != nil {
		return nil, err
	}

	r.cache.Set(ctx, data)

	return data, nil

}
