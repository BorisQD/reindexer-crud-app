package test

import (
	"bytes"
	"context"
	"crud/internal/app"
	"crud/internal/client"
	"crud/internal/domain"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	_ "github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CrudTestSuite struct {
	suite.Suite

	app    *app.App
	client *client.Client

	item domain.Item
}

func (suite *CrudTestSuite) SetupSuite() {
	suite.app = app.New()

	err := suite.app.Config.Load()
	require.NoError(suite.T(), err)
	suite.app.Bootstrap()

	suite.client = client.New(suite.app.Config.DB)

	randomID, _ := uuid.NewV4()
	suite.item = domain.Item{
		Name: "Test name",
		Sort: 1,
		Related: []domain.Nested{
			{
				Name: "Test nested name",
				ID:   randomID,
				Related: []domain.Atom{{
					Name: "Test atom name",
				}},
			},
		},
	}
}

func (suite *CrudTestSuite) SetupTest() {
	ctx := context.Background()
	_ = suite.app.Srv.Start(ctx)
	_ = suite.client.Start(ctx)
}

func (suite *CrudTestSuite) TearDownTest() {
	if suite.item.Empty() {
		return
	}
	_ = suite.client.DeleteItem(context.Background(), suite.item.ID)
}

func (suite *CrudTestSuite) TearDownSuite() {
	suite.app.Shutdown()
}

func (suite *CrudTestSuite) TestAddItemSuccess() {
	item := suite.item
	reqBody, err := json.Marshal(map[string]interface{}{
		"name": item.Name,
		"sort": item.Sort,
		"related": []map[string]string{
			{
				"name": item.Related[0].Name,
				"id":   item.Related[0].ID.String(),
			},
		},
	})
	require.NoError(suite.T(), err)

	resRec := suite.execCreateItemRequest(bytes.NewReader(reqBody), nil)
	assert.Equal(suite.T(), http.StatusOK, resRec.Code)

	var responseBody map[string]string
	err = json.Unmarshal(resRec.Body.Bytes(), &responseBody)
	require.NoError(suite.T(), err)

	id, err := uuid.FromString(responseBody["id"])
	require.NoError(suite.T(), err)

	newItem, ok, err := suite.client.GetItem(context.Background(), id)
	assert.True(suite.T(), ok)

	assert.Equal(suite.T(), item.Name, newItem.Name)
	require.Equal(suite.T(), len(item.Related), len(newItem.Related))
	assert.Equal(suite.T(), item.Related[0].Name, newItem.Related[0].Name)
}

func (suite *CrudTestSuite) execCreateItemRequest(body *bytes.Reader, headers map[string]string) *httptest.ResponseRecorder {
	request, err := http.NewRequest(http.MethodPost, "/items", body)
	require.NoError(suite.T(), err)
	for headerKey, headerValue := range headers {
		request.Header.Set(headerKey, headerValue)
	}
	response := httptest.NewRecorder()

	suite.app.Server.Handler.ServeHTTP(response, request)

	return response
}

//TODO: write other tests

func TestItemCRUDTestSuite(t *testing.T) {
	suite.Run(t, new(CrudTestSuite))
}
