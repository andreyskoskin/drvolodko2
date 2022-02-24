package webapi

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	model "github.com/andreyskoskin/drvolodko2/audit"
)

type auditAPI struct {
	db *postgres
}

func newAuditAPI(db *postgres) *auditAPI {
	return &auditAPI{db}
}

func (api *auditAPI) install(g *echo.Group) {
	g.PUT("/:id", api.put)
}

func (api *auditAPI) put(ctx echo.Context) (err error) {
	var id int64
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return ctx.JSON(http.StatusBadRequest, Message{"invalid id"})
	}

	var audit model.Audit
	if err = ctx.Bind(&audit); err != nil {
		return ctx.JSON(http.StatusBadRequest, Message{err.Error()})
	}

	audit.ID = id
	err = audit.SaveTo(api.db, context.Background())

	if err == model.ErrAuditNotFound {
		return ctx.JSON(http.StatusNotFound, NotFound{
			Entity: "audit",
			ID:     id,
		})
	}

	if err != nil {
		log.Println("ERROR:", err.Error())
		return ctx.JSON(http.StatusInternalServerError, internalError)
	}

	return ctx.JSON(http.StatusOK, updated)
}
