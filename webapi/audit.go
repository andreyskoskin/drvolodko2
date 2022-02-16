package webapi

import (
	"context"
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

func (a *auditAPI) install(g *echo.Group) {
	g.PUT("/:id", a.put)
}

func (a *auditAPI) put(ctx echo.Context) (err error) {
	var id int64
	if id, err = strconv.ParseInt(ctx.Param("id"), 10, 64); err != nil {
		return ctx.JSON(http.StatusBadRequest, BadRequest{Message: "invalid id"})
	}

	var audit model.Audit
	if err = ctx.Bind(&audit); err != nil {
		return ctx.JSON(http.StatusBadRequest, BadRequest{Message: err.Error()})
	}

	audit.ID = id
	err = audit.SaveTo(a.db, context.Background())

	return err
}
