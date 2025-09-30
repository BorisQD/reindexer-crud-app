package client

import (
	"context"
	"crud/internal/domain"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/restream/reindexer"
	_ "github.com/restream/reindexer/v4/bindings/cproto"
)

type Client struct {
	*reindexer.Reindexer
	namespace string
}

func New() *Client {
	db := reindexer.NewReindex("cproto://localhost:6534/reindexer_db", reindexer.WithCreateDBIfMissing())
	return &Client{db, "items422"}
}

func (c Client) Start(ctx context.Context) error {
	clientWithCtx := Client{
		c.WithContext(ctx),
		c.namespace,
	}
	err := clientWithCtx.OpenNamespace(c.namespace, reindexer.DefaultNamespaceOptions(), Item{})
	if err != nil {
		return fmt.Errorf("client.OpenNamespace: %w", err)
	}

	err = clientWithCtx.AddIndex(c.namespace, reindexer.IndexDef{
		Name:      "sorting",
		IndexType: "tree",
		FieldType: "int64",
		JSONPaths: []string{"sort"},
	})
	if err != nil {
		return fmt.Errorf("client.AddIndex: %w", err)
	}

	return nil
}

func (c Client) Stop(ctx context.Context) {
	c.WithContext(ctx).Close()
}

func (c Client) CreateItem(ctx context.Context, item domain.Item) error {
	data := toDTO(item)
	_, err := c.WithContext(ctx).Insert(c.namespace, &data)
	if err != nil {
		return fmt.Errorf("client.CreateItem: %w", err)
	}

	return nil
}

func (c Client) GetItem(ctx context.Context, id uuid.UUID) (domain.Item, bool, error) {
	query := c.WithContext(ctx).Query(c.namespace).Where("id", reindexer.EQ, id.String())

	it := query.Exec()
	defer func() {
		it.Close()
	}()

	var item domain.Item

	if err := it.Error(); err != nil {
		return item, false, fmt.Errorf("client.GetItem: %w", err)
	}

	if it.Count() == 0 {
		return item, false, nil
	}

	it.Next()
	dbItem := it.Object().(*Item)

	return dbItem.toModel(), true, nil
}

func (c Client) GetItems(ctx context.Context, pagination domain.Pagination, order domain.SortOrder) ([]domain.Item, error) {
	query := c.WithContext(ctx).Query(c.namespace).Sort("sort", order == domain.OrderDesc).
		Limit(pagination.Limit).
		Offset(pagination.Offset)

	it := query.Exec()
	defer func() {
		it.Close()
	}()

	if err := it.Error(); err != nil {
		return nil, fmt.Errorf("client.GetItems: %w", err)
	}

	items := make([]domain.Item, 0, it.Count())
	for it.Next() {
		item := it.Object().(*Item)
		items = append(items, item.toModel())
	}

	return items, nil
}

func (c Client) GetItemsCount(ctx context.Context) (int64, error) {
	query := c.WithContext(ctx).Query(c.namespace)

	it := query.Exec()
	defer func() {
		it.Close()
	}()

	if err := it.Error(); err != nil {
		return 0, fmt.Errorf("client.GetItemsCount: %w", err)
	}

	return int64(it.Count()), nil
}

func (c Client) UpdateItem(ctx context.Context, item domain.Item) error {
	dbItem := toDTO(item)
	query := c.WithContext(ctx).Query(c.namespace)

	query.Where("id", reindexer.EQ, dbItem.ID).
		Set("name", dbItem.Name).
		Set("related", dbItem.Related).
		Update()

	if it := query.Exec(); it.Error() != nil {
		return fmt.Errorf("client.UpdateItem: %w", it.Error())
	}

	return nil
}
