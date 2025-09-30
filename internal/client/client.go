package client

import (
	"context"
	"crud/internal/config"
	"crud/internal/domain"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/restream/reindexer"
	_ "github.com/restream/reindexer/v4/bindings/cproto"
	"net"
)

type Client struct {
	*reindexer.Reindexer
	namespace string
}

func New(cfg config.DbConfig) *Client {
	db := reindexer.NewReindex(
		fmt.Sprintf("cproto://%s/%s", net.JoinHostPort(cfg.Host, cfg.Port), cfg.Name),
		reindexer.WithCreateDBIfMissing(),
	)
	return &Client{db, cfg.Namespace}
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
	if !c.IsConnected(ctx) {
		//закрытие закрытого канала под капотом Close вызовет панику
		return
	}
	_ = c.CloseNamespace(c.namespace)
	c.WithContext(ctx).Close()
}

func (c Client) IsConnected(ctx context.Context) bool {
	err := c.WithContext(ctx).Ping()
	return err == nil
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

	var item domain.Item

	it := query.Exec()
	if err := it.Error(); err != nil {
		return item, false, fmt.Errorf("client.GetItem: %w", err)
	}

	defer func() {
		it.Close()
	}()

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
	if err := it.Error(); err != nil {
		return nil, fmt.Errorf("client.GetItems: %w", err)
	}

	defer func() {
		it.Close()
	}()

	items := make([]domain.Item, 0, it.Count())
	for it.Next() {
		item := it.Object().(*Item)
		items = append(items, item.toModel())
	}

	return items, nil
}

func (c Client) GetItemsCount(ctx context.Context) (int64, error) {
	query := c.WithContext(ctx).Query(c.namespace).ReqTotal()

	it := query.Exec()
	if err := it.Error(); err != nil {
		return 0, fmt.Errorf("client.GetItemsCount: %w", err)
	}

	defer func() {
		it.Close()
	}()

	res := it.AggResults()[0]

	return int64(res.Value), nil
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

func (c Client) DeleteItem(ctx context.Context, id uuid.UUID) error {
	err := c.WithContext(ctx).Delete(c.namespace, &Item{ID: id.String()})
	if err != nil {
		return fmt.Errorf("client.DeleteItem: %w", err)
	}

	return nil
}
