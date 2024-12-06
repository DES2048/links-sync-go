package webapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"links-sync-go/internal/storage"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

func getIdParam(c echo.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return id, echo.NewHTTPError(400, "invalid id")
	}
	return id, nil
}

func readBodyJson(r *http.Request, value any) error {
	if err := json.NewDecoder(r.Body).Decode(value); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to parse json")
	}
	return nil
}

type VisitedHandler struct {
	server IApiServer
}

func NewVisitedHandler(server IApiServer) *VisitedHandler {
	return &VisitedHandler{
		server: server,
	}
}

func (handler *VisitedHandler) saveVisitedHandler(c echo.Context) error {
	visitedData := []*storage.DbCreateVisited{}

	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(400, err.Error())
	}

	err = json.Unmarshal(data, &visitedData)
	if err != nil {
		return c.String(400, err.Error())
	}

	err = validation.Validate(visitedData)
	if err != nil {
		return c.JSONPretty(400, err, " ")
	}

	// set default status value
	for _, vis := range visitedData {
		if vis.Status == 0 {
			vis.Status = 1
		}
	}

	repo := handler.server.Storage().VisitedRepo()
	err = repo.AddBatch(visitedData)
	if err != nil {
		return c.String(500, fmt.Sprintf("%s", err))
	}
	return c.JSONPretty(200, "ok", "  ")
}

func (handler *VisitedHandler) getVisitedByIds(c echo.Context) error {
	idsData := visitedIdsPayload{}

	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(400, err.Error())
	}

	err = json.Unmarshal(data, &idsData)
	if err != nil {
		return c.String(400, err.Error())
	}

	visited, err := handler.server.Storage().VisitedRepo().ListByIds(idsData.Ids)
	if err != nil {
		return c.String(500, err.Error())
	}

	resp := make([]*IdStatusReponse, 0, len(visited))

	for _, v := range visited {
		resp = append(resp, &IdStatusReponse{
			Id:     v.Id,
			Status: v.Status,
		})
	}
	return c.JSONPretty(200, resp, " ")
}

func (handler *VisitedHandler) getVisitedListHandler(c echo.Context) error {
	repo := handler.server.Storage().VisitedRepo()

	data, err := repo.List()
	if err != nil {
		return c.String(500, err.Error())
	}

	return c.JSONPretty(200, data, " ")
}

func (handler *VisitedHandler) getVisitedHandler(c echo.Context) error {
	id, err := getIdParam(c)
	if err != nil {
		return err
	}
	v, err := handler.server.Storage().VisitedRepo().Get(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return echo.NewHTTPError(400, err)
		}
		return err
	}

	return c.JSONPretty(200, v, " ")
}

func (handler *VisitedHandler) deleteVisitedHandler(c echo.Context) error {
	id, err := getIdParam(c)
	if err != nil {
		return err
	}

	err = handler.server.Storage().VisitedRepo().Delete(id)
	if err != nil {
		return c.String(400, err.Error())
	}

	return c.NoContent(204)
}

func (handler *VisitedHandler) patchVisitedHandler2(c echo.Context) error {
	id, err := getIdParam(c)
	if err != nil {
		return err
	}

	patchData := new(storage.DbPatchVisited)

	err = readBodyJson(c.Request(), patchData)
	if err != nil {
		return err
	}

	// validate
	err = validation.Validate(patchData)
	if err != nil {
		return echo.NewHTTPError(400, err)
	}

	err = handler.server.Storage().VisitedRepo().UpdatePartial(id, patchData)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return echo.NewHTTPError(400, err)
		}
		return err
	}

	return c.String(200, "patch visited")
}
