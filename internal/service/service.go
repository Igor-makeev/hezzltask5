package service

import (
	"context"
	"hezzltask5/internal/models"
	"hezzltask5/internal/repository"
	"hezzltask5/internal/service/queue"
	"time"
)

type Service struct {
	Repo  *repository.Repository
	Queue *queue.NatsQueue
}

func NewService(r *repository.Repository, q *queue.NatsQueue) *Service {
	return &Service{Repo: r, Queue: q}
}

func (s *Service) Create(ctx context.Context, campaignId int, name string) (models.Item, error) {
	item, err := s.Repo.Create(ctx, campaignId, name)
	s.Queue.AddToQueue(models.Event{
		Id:          item.Id,
		CampaignId:  item.CampaignId,
		Name:        item.Name,
		Description: item.Description,
		Priority:    item.Priority,
		Removed:     item.Removed,
		EventTime:   item.CreatedAt,
	})
	return item, err
}

func (s *Service) Update(ctx context.Context, id, campaignId int, name, description string) (models.Item, error) {
	item, err := s.Repo.Update(ctx, id, campaignId, name, description)
	s.Queue.AddToQueue(models.Event{
		Id:          item.Id,
		CampaignId:  item.CampaignId,
		Name:        item.Name,
		Description: item.Description,
		Priority:    item.Priority,
		Removed:     item.Removed,
		EventTime:   time.Now(),
	})
	return item, err
}

func (s *Service) Delete(ctx context.Context, id, campaignId int) (models.DeleteResponse, error) {
	item, err := s.Repo.Delete(ctx, id, campaignId)
	if err != nil {
		return models.DeleteResponse{}, err
	}
	s.Queue.AddToQueue(models.Event{
		Id:          item.Id,
		CampaignId:  item.CampaignId,
		Name:        item.Name,
		Description: item.Description,
		Priority:    item.Priority,
		Removed:     item.Removed,
		EventTime:   time.Now(),
	})
	response := models.DeleteResponse{
		Id:         item.Id,
		CampaignId: item.CampaignId,
		Removed:    item.Removed,
	}
	return response, nil
}
func (s *Service) GetList(ctx context.Context) ([]models.Item, error) {
	return s.Repo.GetList(ctx)
}
