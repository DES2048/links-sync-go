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

func TestVisitedHandlers_SaveVisited(t *testing.T) {

	server := TestServer(t)

	// test happy case
	data := []*storage.DbCreateVisited{
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

	// create want
	want := make([]*storage.DbVisited, 0, len(data))

	for _, e := range data {
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

	dataJson, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Error while marshall to json %s\n", err)
	}

	e := echo.New()
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(dataJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//req.Header.Set(echo.HeaderAuthorization, "Bearer ")
	resp := httptest.NewRecorder()

	c := e.NewContext(req, resp)

	handler := NewVisitedHandler(server)

	err = handler.saveVisitedHandler(c)

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)

	/*
		var respBoby any
		_ = json.NewDecoder(resp.Body).Decode(&respBoby)
		t.Logf("body is %#v\n", respBoby)
	*/

	// check results
	ids := utils.SliceMap(data, func(e *storage.DbCreateVisited) int {
		return e.Id
	})
	gotVisited, err := server.Storage().VisitedRepo().ListByIds(ids)

	assert.Nil(t, err)
	assert.Equal(t, len(data), len(gotVisited))

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

	// test empty array
	testData := []*storage.DbCreateVisited{}
	dataJson, _ = json.Marshal(testData)
	req, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader(dataJson))
	resp = httptest.NewRecorder()

	c = e.NewContext(req, resp)
	err = handler.saveVisitedHandler(c)

	assert.NotNil(t, err)

	httpErr := &echo.HTTPError{}

	assert.ErrorAs(t, err, &httpErr)

	assert.Equal(t, httpErr.Code, 400)

	// test required
	testData = []*storage.DbCreateVisited{
		{},
	}
	dataJson, _ = json.Marshal(testData)
	req, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader(dataJson))
	resp = httptest.NewRecorder()

	c = e.NewContext(req, resp)
	err = handler.saveVisitedHandler(c)

	assert.NotNil(t, err)

	//errs := make(validation.Errors)

	assert.ErrorAs(t, err, &httpErr)

	assert.Equal(t, httpErr.Code, 400)

	// test status
}
