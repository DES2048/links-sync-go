package webapi

import (
	"bytes"
	"encoding/json"
	"links-sync-go/internal/storage"
	"links-sync-go/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var visitedDataCreate = []*storage.DbCreateVisited{
	{
		Id:        1,
		Title:     "title1",
		PosterUrl: "https://some_url",
	},
	{
		Id:        2,
		Title:     "title2",
		PosterUrl: "https://some_url2",
		Status:    2,
	},
	{
		Id:    3,
		Title: "title3",
	},
}

func createRequest(method string, url string, body interface{}) *http.Request {
	dataJson, err := json.Marshal(body)
	if err != nil {
		panic("Error while marshall to json")
	}

	req, _ := http.NewRequest(method, url, bytes.NewReader(dataJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return req
}

func TestVisitedHandlers_SaveVisited(t *testing.T) {
	server := TestServer(t)

	// TEST happy case

	// create want
	want := make([]*storage.DbVisited, 0, len(visitedDataCreate))

	for _, e := range visitedDataCreate {
		status := e.Status
		if status == 0 {
			status = 1 // default value
		}

		want = append(want, &storage.DbVisited{
			Id:        e.Id,
			Title:     e.Title,
			PosterUrl: e.PosterUrl,
			Status:    status,
		})
	}

	e := echo.New()
	req := createRequest(http.MethodPost, "/", visitedDataCreate)
	resp := httptest.NewRecorder()

	c := e.NewContext(req, resp)

	handler := NewVisitedHandler(server)

	err := handler.saveVisitedHandler(c)

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)

	/*
		var respBoby any
		_ = json.NewDecoder(resp.Body).Decode(&respBoby)
		t.Logf("body is %#v\n", respBoby)
	*/

	// check results
	ids := utils.SliceMap(visitedDataCreate, func(e *storage.DbCreateVisited) int {
		return e.Id
	})
	gotVisited, err := server.Storage().VisitedRepo().ListByIds(ids)

	assert.Nil(t, err)
	assert.Equal(t, len(visitedDataCreate), len(gotVisited))

	// drop add date
	for _, e := range gotVisited {
		e.AddDate = ""
	}

	/*
		slices.SortFunc(gotVisited, func(a *storage.DbVisited, b *storage.DbVisited) int {
			return cmp.Compare(a.Id, b.Id)
		}) */
	assert.ElementsMatch(t, want, gotVisited)

	// test validation

	// TEST empty array
	testData := []*storage.DbCreateVisited{}
	req = createRequest(http.MethodPost, "/", testData)
	resp = httptest.NewRecorder()

	c = e.NewContext(req, resp)
	err = handler.saveVisitedHandler(c)

	assert.NotNil(t, err)

	httpErr := &echo.HTTPError{}

	assert.ErrorAs(t, err, &httpErr)

	assert.Equal(t, 400, httpErr.Code)

	// test required
	testData = []*storage.DbCreateVisited{
		{},
	}
	req = createRequest(http.MethodPost, "/", testData)
	resp = httptest.NewRecorder()

	c = e.NewContext(req, resp)
	err = handler.saveVisitedHandler(c)

	assert.NotNil(t, err)

	// errs := make(validation.Errors)

	assert.ErrorAs(t, err, &httpErr)

	assert.Equal(t, 400, httpErr.Code)

	// test status
}

func TestVisitedHandlers_GetByIds(t *testing.T) {
	server := TestServer(t)

	data := visitedIdsPayload{
		[]int{1, 2, 4},
	}

	server.Storage().VisitedRepo().AddBatch(visitedDataCreate)

	e := echo.New()

	req := createRequest(http.MethodPost, "/", data)
	resp := httptest.NewRecorder()

	c := e.NewContext(req, resp)

	handler := NewVisitedHandler(server)

	err := handler.getVisitedByIds(c)

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)

	var respBoby []IdStatusReponse
	_ = json.NewDecoder(resp.Body).Decode(&respBoby)
	expected := []IdStatusReponse{
		{
			visitedDataCreate[0].Id,
			1,
		},
		{
			visitedDataCreate[1].Id,
			visitedDataCreate[1].Status,
		},
	}

	assert.ElementsMatch(t, expected, respBoby)
}
