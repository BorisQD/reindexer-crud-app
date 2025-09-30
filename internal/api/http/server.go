package http

import (
	"context"
	"crud/internal/domain"
	"errors"
	"github.com/gofrs/uuid/v5"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"log/slog"
)

type service interface {
	CreateItem(ctx context.Context, item domain.Item) (uuid.UUID, error)
	GetItem(ctx context.Context, id uuid.UUID) (domain.Item, bool, error)
	GetItemsPaginated(ctx context.Context, pagination domain.Pagination) ([]domain.Item, int64, error)
	UpdateItem(ctx context.Context, id uuid.UUID, item domain.Item) error
}

type healthChecker interface {
	HealthCheck(ctx context.Context) bool
}

type Server struct {
	Service service
	Checker healthChecker
	Logger  *slog.Logger
}

func (s Server) GetHealth(ctx context.Context, _ GetHealthRequestObject) (GetHealthResponseObject, error) {
	if ok := s.Checker.HealthCheck(ctx); !ok {
		return GetHealth503JSONResponse{}, errors.New("not healthy")
	}

	return GetHealth200JSONResponse{}, nil
}

func (s Server) GetLive(_ context.Context, _ GetLiveRequestObject) (GetLiveResponseObject, error) {
	return GetLive200JSONResponse{}, nil
}

func (s Server) GetItems(ctx context.Context, request GetItemsRequestObject) (GetItemsResponseObject, error) {
	items, totalCount, err := s.Service.GetItemsPaginated(ctx, domain.Pagination{
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
			Id:        openapitypes.UUID(it.ID),
			Name:      it.Name,
			UpdatedAt: it.UpdatedAt,
		})
	}

	total := int(totalCount)
	return GetItems200JSONResponse{
		Items: res,
		Total: &total,
	}, nil
}

func (s Server) PostItems(ctx context.Context, request PostItemsRequestObject) (PostItemsResponseObject, error) {
	itemID, err := s.Service.CreateItem(ctx, createRequestToItem(*request.Body))

	if err != nil {
		s.Logger.Error(err.Error())
		return PostItems500JSONResponse{}, err
	}

	return PostItems200JSONResponse{
		Id: openapitypes.UUID(itemID),
	}, nil
}

func (s Server) GetItemsId(ctx context.Context, request GetItemsIdRequestObject) (GetItemsIdResponseObject, error) {
	item, ok, err := s.Service.GetItem(ctx, uuid.UUID(request.Id))
	if err != nil {
		s.Logger.Error(err.Error())
		return GetItemsId500JSONResponse{}, err
	}
	if !ok {
		return GetItemsId404JSONResponse{}, err
	}

	nested := make([]Nested, 0, len(item.Related))
	for _, nst := range item.Related {
		nested = append(nested, Nested{
			Id:   openapitypes.UUID(nst.ID),
			Name: nst.Name,
		})
	}
	return GetItemsId200JSONResponse{
		Name:      item.Name,
		CreatedAt: item.CreatedAt,
		Related:   &nested,
	}, nil
}

func (s Server) PutItemsId(ctx context.Context, request PutItemsIdRequestObject) (PutItemsIdResponseObject, error) {
	err := s.Service.UpdateItem(ctx, uuid.UUID(request.Id), updateRequestToItem(*request.Body))
	if err != nil {
		s.Logger.Error(err.Error())
		return PutItemsId500JSONResponse{}, err
	}

	return PutItemsId200JSONResponse{}, nil
}

func createRequestToItem(req ItemCreate) domain.Item {
	nested := make([]domain.Nested, 0, len(*req.Related))
	for _, nst := range nested {
		nested = append(nested, domain.Nested{
			ID:   nst.ID,
			Name: nst.Name,
		})
	}

	return domain.Item{
		Name:    req.Name,
		Related: nested,
	}
}

func updateRequestToItem(req ItemUpdate) domain.Item {
	nested := make([]domain.Nested, 0, len(req.Related))
	for _, nst := range nested {
		nested = append(nested, domain.Nested{
			ID:   nst.ID,
			Name: nst.Name,
		})
	}

	return domain.Item{
		Name:    req.Name,
		Related: nested,
	}
}
