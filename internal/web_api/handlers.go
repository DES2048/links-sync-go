package webapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"links-sync-go/internal/storage"
	"net/http"
	"strconv"

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
	server      IApiServer
	visitedRepo storage.IVisitedRepo
}

func NewVisitedHandler(server IApiServer) *VisitedHandler {
	return &VisitedHandler{
		server:      server,
		visitedRepo: server.Storage().VisitedRepo(),
	}
}

func (handler *VisitedHandler) saveVisitedHandler(c echo.Context) error {
	visitedData := []*storage.DbCreateVisited{}

	err := c.Bind(&visitedData)
	if err != nil {
		return err
	}

	if len(visitedData) == 0 {
		return echo.NewHTTPError(400, "Empty array passed")
	}

	err = validation.Validate(visitedData)
	if err != nil {
		return echo.NewHTTPError(400, err)
	}

	err = handler.visitedRepo.AddBatch(visitedData)
	if err != nil {
		return c.String(500, fmt.Sprintf("%s", err))
	}
	return c.JSONPretty(200, "ok", "  ")
}

func (handler *VisitedHandler) getVisitedByIds(c echo.Context) error {
	idsData := visitedIdsPayload{}
	err := c.Bind(&idsData)
	if err != nil {
		return err
	}
	visited, err := handler.visitedRepo.ListByIds(idsData.Ids)
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
	data, err := handler.visitedRepo.List()
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
	v, err := handler.visitedRepo.Get(id)
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

	err = handler.visitedRepo.Delete(id)
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

	err = handler.visitedRepo.UpdatePartial(id, patchData)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return echo.NewHTTPError(400, err)
		}
		return err
	}

	return c.String(200, "ok")
}
