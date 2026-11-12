package webapi

import (
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

func saveVisitedHandler(visitedRepo storage.IVisitedRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
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
			return err
		}

		err = visitedRepo.AddBatch(visitedData)
		if err != nil {
			return err
		}
		return c.JSONPretty(200, "ok", "  ")
	}
}

func getVisitedByIds(visitedRepo storage.IVisitedRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		idsData := visitedIdsPayload{}
		err := c.Bind(&idsData)
		if err != nil {
			return err
		}
		visited, err := visitedRepo.ListByIds(idsData.Ids)
		if err != nil {
			return err
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
}

func getVisitedListHandler(visitedRepo storage.IVisitedRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		limit := new(30)
		offset := new(0)
		err := echo.QueryParamsBinder(c).
			Int("start", offset).
			Int("items", limit).
			BindError()
		if err != nil {
			return err
		}

		if *offset > 0 {
			*offset--
		}
		data, err := visitedRepo.List(*limit, *offset)
		if err != nil {
			return err
		}

		return c.JSONPretty(200, data, " ")
	}
}

func getVisitedHandler(visitedRepo storage.IVisitedRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := new(int)
		err := echo.PathParamsBinder(c).Int("id", id).BindError()
		if err != nil {
			return err
		}
		v, err := visitedRepo.Get(*id)
		if err != nil {
			return err
		}

		return c.JSONPretty(200, v, " ")
	}
}

func deleteVisitedHandler(visitedRepo storage.IVisitedRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := getIdParam(c)
		if err != nil {
			return err
		}

		err = visitedRepo.Delete(id)
		if err != nil {
			return err
		}

		return c.NoContent(204)
	}
}

func patchVisitedHandler2(visitedRepo storage.IVisitedRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := getIdParam(c)
		if err != nil {
			return err
		}

		patchData := new(storage.DbPatchVisited)

		err = c.Bind(patchData)
		if err != nil {
			return err
		}

		// validate
		err = validation.Validate(patchData)
		if err != nil {
			return err
		}

		err = visitedRepo.UpdatePartial(id, patchData)
		if err != nil {
			return err
		}

		return c.String(200, "ok")
	}
}
