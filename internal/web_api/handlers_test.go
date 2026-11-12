package webapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"links-sync-go/internal/storage"
	"links-sync-go/internal/utils"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var visitedDataCreate = []*storage.DbCreateVisited{
	{
		Id:        1,
		Title:     "title1",
		PosterUrl: new("https://some_url"),
	},
	{
		Id:        2,
		Title:     "title2",
		PosterUrl: new("https://some_url2"),
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

	e := server.Echo()
	req := createRequest(http.MethodPost, "/", visitedDataCreate)
	resp := httptest.NewRecorder()

	c := e.NewContext(req, resp)

	handler := saveVisitedHandler(server.Storage().VisitedRepo())

	err := handler(c)

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

	t.Run("empty visited", func(t *testing.T) {
		// TEST empty array
		testData := []*storage.DbCreateVisited{}
		req = createRequest(http.MethodPost, "/", testData)
		resp = httptest.NewRecorder()

		c = e.NewContext(req, resp)
		err = handler(c)

		assert.NotNil(t, err)

		httpErr := &echo.HTTPError{}

		assert.ErrorAs(t, err, &httpErr)

		assert.Equal(t, 400, httpErr.Code)
	})

	t.Run("required fields", func(t *testing.T) {
		testData := []*storage.DbCreateVisited{
			{},
		}
		req = createRequest(http.MethodPost, "/save-links", testData)
		resp = httptest.NewRecorder()

		c = e.NewContext(req, resp)
		err = handler(c)

		assert.NotNil(t, err)

		// errs := make(validation.Errors)
		var validErr validation.Errors
		assert.ErrorAs(t, err, &validErr)
	})
}

func TestVisitedHandlers_GetByIds(t *testing.T) {
	server := TestServer(t)

	data := visitedIdsPayload{
		[]int{1, 2, 4},
	}

	server.Storage().VisitedRepo().AddBatch(visitedDataCreate)

	e := server.Echo()

	req := createRequest(http.MethodPost, "/", data)
	resp := httptest.NewRecorder()

	c := e.NewContext(req, resp)

	handler := getVisitedByIds(server.Storage().VisitedRepo())

	err := handler(c)

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

func TestGetVisitedHandler(t *testing.T) {
	server := TestServer(t)

	e := server.Echo()
	handler := getVisitedHandler(server.Storage().VisitedRepo())

	visited := &storage.DbCreateVisited{
		Id:     1,
		Title:  "title",
		Status: 1,
	}

	server.Storage().VisitedRepo().Add(visited)

	t.Run("not found", func(t *testing.T) {
		req := createRequest(http.MethodGet, "/get-links/2", nil)
		resp := httptest.NewRecorder()

		c := e.NewContext(req, resp)
		c.SetPath("/get-links/:id")
		c.SetParamNames("id")
		c.SetParamValues("2")

		err := handler(c)
		assert.True(t, storage.IsErrNotFound(err))
	})

	t.Run("found", func(t *testing.T) {
		req := createRequest(http.MethodGet, "/get-links/1", nil)
		resp := httptest.NewRecorder()

		c := e.NewContext(req, resp)
		c.SetPath("/get-links/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		err := handler(c)
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.Code)
		var respBoby *storage.DbVisited
		err = json.NewDecoder(resp.Body).Decode(&respBoby)
		assert.Nil(t, err)
		respBoby.AddDate = ""
		assert.EqualValues(t, &storage.DbVisited{
			Id:     visited.Id,
			Title:  visited.Title,
			Status: visited.Status,
		}, respBoby)
	})
}

func TestDeleteVisitedHandler(t *testing.T) {
	server := TestServer(t)

	e := server.Echo()
	handler := deleteVisitedHandler(server.Storage().VisitedRepo())

	visited := &storage.DbCreateVisited{
		Id:     1,
		Title:  "title",
		Status: 1,
	}

	server.Storage().VisitedRepo().Add(visited)

	t.Run("not found", func(t *testing.T) {
		req := createRequest(http.MethodDelete, "/get-links/2", nil)
		resp := httptest.NewRecorder()

		c := e.NewContext(req, resp)
		c.SetPath("/get-links/:id")
		c.SetParamNames("id")
		c.SetParamValues("2")

		err := handler(c)
		assert.True(t, storage.IsErrNotFound(err))
	})

	t.Run("del", func(t *testing.T) {
		req := createRequest(http.MethodDelete, "/get-links/1", nil)
		resp := httptest.NewRecorder()

		c := e.NewContext(req, resp)
		c.SetPath("/get-links/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		err := handler(c)
		assert.Nil(t, err)
		assert.Equal(t, 204, resp.Code)
	})
}

func TestPatchVisitedHandler(t *testing.T) {
	server := TestServer(t)

	e := server.Echo()
	handler := patchVisitedHandler2(server.Storage().VisitedRepo())

	visited := &storage.DbCreateVisited{
		Id:     1,
		Title:  "title",
		Status: 1,
	}

	server.Storage().VisitedRepo().Add(visited)

	t.Run("not found", func(t *testing.T) {
		req := createRequest(http.MethodPatch, "/get-links/2", map[string]any{"status": 2})
		resp := httptest.NewRecorder()

		c := e.NewContext(req, resp)
		c.SetPath("/get-links/:id")
		c.SetParamNames("id")
		c.SetParamValues("2")

		err := handler(c)
		assert.True(t, storage.IsErrNotFound(err))
	})

	t.Run("status", func(t *testing.T) {
		req := createRequest(http.MethodPatch, "/get-links/1", map[string]any{"status": 2})
		resp := httptest.NewRecorder()

		c := e.NewContext(req, resp)

		c.SetPath("/get-links/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		err := handler(c)
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.Code)

		// get patched
		patched, _ := server.Storage().VisitedRepo().Get(1)

		assert.Equal(t, 2, patched.Status)
	})
}
