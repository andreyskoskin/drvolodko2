package auditor

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/andreyskoskin/drvolodko2/api"
	"github.com/andreyskoskin/drvolodko2/datamodel"
)

type echoAPI struct {
	auditors datamodel.Auditors
}

func NewEchoAPI(auditors datamodel.Auditors) api.EchoAPI {
	return &echoAPI{
		auditors: auditors,
	}
}

func (api *echoAPI) Bind(g *echo.Group) {
	g.GET("/:id", api.GetOne)
	g.GET("/", api.GetMany)
	g.POST("/", api.PostOne)
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

	var reply *datamodel.Auditor
	reply, err = api.auditors.FindOne(datamodel.AuditorID(request.ID))

	if err == datamodel.ErrAuditorNotFound {
		return ctx.String(http.StatusNotFound, err.Error())
	}

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, reply)
}

// ----------------------------------------------------------------

type GetManyReply struct {
	Items []datamodel.Auditor `json:"items"`
}

func (api *echoAPI) GetMany(ctx echo.Context) (err error) {
	var items []datamodel.Auditor
	items, err = api.auditors.FindAll()

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, GetManyReply{
		Items: items,
	})
}

// ----------------------------------------------------------------

type PostOneRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	PositionID int64  `json:"position_id"`
}

type PostOneReply struct {
	ID int `param:"id"`
}

func (api *echoAPI) PostOne(ctx echo.Context) (err error) {
	var request PostOneRequest
	err = ctx.Bind(&request)

	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	if request.FirstName == "" {
		return ctx.String(http.StatusBadRequest, "empty first name")
	}

	if request.LastName == "" {
		return ctx.String(http.StatusBadRequest, "empty last name")
	}

	var id datamodel.AuditorID
	id, err = api.auditors.Create(datamodel.Auditor{
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		PositionID: request.PositionID,
	})

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, PostOneReply{
		ID: int(id),
	})
}
