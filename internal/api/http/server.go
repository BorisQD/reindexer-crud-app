package http

import (
	"context"
	"crud/internal/domain"
	"github.com/gofrs/uuid/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type service interface {
	CreateItem(ctx context.Context, item domain.Item) error
	GetItem(ctx context.Context, id uuid.UUID) (domain.Item, bool, error)
	GetItemsPaginated(ctx context.Context, pagination domain.Pagination) ([]domain.Item, error)
	UpdateItem(ctx context.Context, id uuid.UUID, item domain.Item) error
}

type Server struct {
	service service
}

func (s Server) GetHealth(ctx context.Context, request GetHealthRequestObject) (GetHealthResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetItems(ctx context.Context, request GetItemsRequestObject) (GetItemsResponseObject, error) {
	items, err := s.service.GetItemsPaginated(ctx, domain.Pagination{
		Limit:  request.Params.Limit,
		Offset: request.Params.Offset,
	})

	if err != nil {
		return GetItems500JSONResponse{}, err
	}

	res := make([]Item, 0, len(items))
	for _, it := range items {
		res = append(res, Item{
			CreatedAt: it.CreatedAt,
			Id:        openapi_types.UUID(it.ID),
			Name:      it.Name,
			UpdatedAt: it.UpdatedAt,
		})
	}

	return GetItems200JSONResponse{
		Items: res,
	}, nil
}

func (s Server) PostItems(ctx context.Context, request PostItemsRequestObject) (PostItemsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetItemsId(ctx context.Context, request GetItemsIdRequestObject) (GetItemsIdResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) PutItemsId(ctx context.Context, request PutItemsIdRequestObject) (PutItemsIdResponseObject, error) {
	//TODO implement me
	panic("implement me")
}
