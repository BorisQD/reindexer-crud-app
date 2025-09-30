package service

import (
	"context"
	"crud/internal/domain"
	"github.com/gofrs/uuid/v5"
	"github.com/jellydator/ttlcache/v3"
	"slices"
	"sync"
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
	db    dbClient
	cache *ttlcache.Cache[uuid.UUID, domain.Item]
}

func New(db dbClient, ttl time.Duration) *Service {
	itemCache := ttlcache.New(
		ttlcache.WithDisableTouchOnHit[uuid.UUID, domain.Item](),
		ttlcache.WithTTL[uuid.UUID, domain.Item](ttl),
	)
	return &Service{db, itemCache}
}

func (s Service) Start(ctx context.Context) error {
	go s.cache.Start()
	return s.db.Start(ctx)
}

func (s Service) Close(ctx context.Context) {
	go s.cache.Stop()
	s.db.Stop(ctx)
}

func (s Service) CreateItem(ctx context.Context, item domain.Item) (newItemID uuid.UUID, err error) {
	if item.Empty() {
		id, _ := uuid.NewV4()
		item.ID = id
	}
	item.CreatedAt = time.Now()
	item.UpdatedAt = &item.CreatedAt

	return item.ID, s.db.CreateItem(ctx, item)
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

	return processItems(items, func(item domain.Item) domain.Item {
		it := domain.Transform(item)
		slices.SortFunc(it.Related, func(a, b domain.Nested) int {
			return int(b.Sort - a.Sort)
		})

		return it
	}), total, nil
}

func (s Service) GetItem(ctx context.Context, id uuid.UUID) (domain.Item, bool, error) {
	cached := s.cache.Get(id)
	if cached != nil && !cached.IsExpired() {
		return cached.Value(), true, nil
	}

	item, found, err := s.db.GetItem(ctx, id)
	if err == nil {
		s.cache.Set(id, item, ttlcache.DefaultTTL)
	}

	return item, found, err
}

func (s Service) UpdateItem(ctx context.Context, id uuid.UUID, item domain.Item) error {
	item.ID = id
	return s.db.UpdateItem(ctx, item)
}

func processItems(items []domain.Item, transform func(item domain.Item) domain.Item) []domain.Item {
	ch := make(chan struct {
		index int
		value domain.Item
	}, len(items))

	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		go func(idx int, it domain.Item) {
			defer wg.Done()

			transformed := transform(it)
			ch <- struct {
				index int
				value domain.Item
			}{idx, transformed}
		}(i, item)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		items[result.index] = result.value
	}

	return items
}
