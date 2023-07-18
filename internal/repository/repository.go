package repository

import (
	"context"
	"hezzltask5/internal/models"

	"github.com/sirupsen/logrus"
)

type persistanceStorage interface {
	Create(ctx context.Context, campaignId int, name string) (models.Item, error)
	Update(ctx context.Context, id, campaignId int, name, description string) (models.Item, error)
	Delete(ctx context.Context, id, campaignId int) (models.Item, error)
	GetList(ctx context.Context) ([]models.Item, error)
}

type cache interface {
	Get(ctx context.Context) ([]models.Item, error)
	Set(ctx context.Context, items []models.Item)
	Remove(ctx context.Context)
}

type Repository struct {
	persistanceStorage persistanceStorage
	cache              cache
}

func NewRepository(ps persistanceStorage, c cache) *Repository {
	return &Repository{persistanceStorage: ps, cache: c}
}

func (r *Repository) Create(ctx context.Context, campaignId int, name string) (models.Item, error) {
	return r.persistanceStorage.Create(ctx, campaignId, name)
}

func (r *Repository) Update(ctx context.Context, id, campaignId int, name, description string) (models.Item, error) {
	r.cache.Remove(ctx)
	return r.persistanceStorage.Update(ctx, id, campaignId, name, description)
}

func (r *Repository) Delete(ctx context.Context, id, campaignId int) (models.Item, error) {
	r.cache.Remove(ctx)
	return r.persistanceStorage.Delete(ctx, id, campaignId)
}

func (r *Repository) GetList(ctx context.Context) ([]models.Item, error) {

	if data, err := r.cache.Get(ctx); data != nil {
		return data, nil
	} else {
		logrus.Printf("repo level: getlist method: %v", err)
	}
	data, err := r.persistanceStorage.GetList(ctx)
	if err != nil {
		return nil, err
	}

	r.cache.Set(ctx, data)

	return data, nil

}
