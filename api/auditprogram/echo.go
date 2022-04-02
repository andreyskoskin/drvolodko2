package auditprogram

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/andreyskoskin/drvolodko2/api"
	"github.com/andreyskoskin/drvolodko2/datamodel"
)

type echoAPI struct {
	auditPrograms datamodel.AuditPrograms
}

func NewEchoAPI(auditPrograms datamodel.AuditPrograms) api.EchoAPI {
	return &echoAPI{
		auditPrograms: auditPrograms,
	}
}

func (api *echoAPI) Bind(g *echo.Group) {
	g.GET("/:id", api.GetOne)
	g.PUT("/:id", api.PutOne)
}

// ----------------------------------------------------------------

type GetOneRequest struct {
	ID int `param:"id"`
}

func (api *echoAPI) GetOne(ctx echo.Context) (err error) {
	var request GetOneRequest
	err = ctx.Bind(&request)

	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	var reply *datamodel.AuditProgram
	reply, err = api.auditPrograms.FindOne(datamodel.AuditProgramID(request.ID))

	if err == datamodel.ErrAuditProgramNotFound {
		return ctx.String(http.StatusNotFound, err.Error())
	}

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, reply)
}

// ----------------------------------------------------------------

type PutOneRequest struct {
	datamodel.AuditProgram
	ID int `param:"id"`
}

func (api *echoAPI) PutOne(ctx echo.Context) (err error) {
	var request PutOneRequest
	err = ctx.Bind(&request)

	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	err = api.auditPrograms.Update(request.AuditProgram)
	if err == datamodel.ErrAuditProgramNotFound {
		return ctx.String(http.StatusNotFound, err.Error())
	}

	if err != nil {
		return err
	}

	return ctx.String(http.StatusOK, "updated")
}
