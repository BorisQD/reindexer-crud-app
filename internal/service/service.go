package service

import (
	"context"
	"crud/internal/domain"
	"github.com/gofrs/uuid/v5"
	"time"
)

type dbClient interface {
	CreateItem(ctx context.Context, item domain.Item) error
	GetItem(ctx context.Context, id uuid.UUID) (domain.Item, bool, error)
	GetItems(ctx context.Context, pagination domain.Pagination, order domain.SortOrder) ([]domain.Item, error)
	GetItemsCount(ctx context.Context) (int64, error)
	UpdateItem(ctx context.Context, item domain.Item) error

	Start(ctx context.Context) error
	Stop(ctx context.Context)
}

type Service struct {
	db dbClient
}

func New(db dbClient) *Service {
	return &Service{db}
}

func (s Service) Start(ctx context.Context) error {
	return s.db.Start(ctx)
}

func (s Service) Close(ctx context.Context) {
	s.db.Stop(ctx)
}

func (s Service) CreateItem(ctx context.Context, item domain.Item) error {
	if item.ID == uuid.Nil {
		id, _ := uuid.NewV4()
		item.ID = id
	}
	item.CreatedAt = time.Now()
	item.UpdatedAt = &item.CreatedAt

	return s.db.CreateItem(ctx, item)
}

func (s Service) GetItemsPaginated(ctx context.Context, pagination domain.Pagination) (items []domain.Item, total int64, err error) {
	items, err = s.db.GetItems(ctx, pagination, domain.OrderDesc)
	if err != nil {
		return items, total, err
	}

	total, err = s.db.GetItemsCount(ctx)
	if err != nil {
		return items, total, err
	}

	return items, total, nil
}

func (s Service) GetItem(ctx context.Context, id uuid.UUID) (domain.Item, bool, error) {
	return s.db.GetItem(ctx, id)
}

func (s Service) UpdateItem(ctx context.Context, id uuid.UUID, item domain.Item) error {
	item.ID = id
	return s.db.UpdateItem(ctx, item)
}
