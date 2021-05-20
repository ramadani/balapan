package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/ramadani/balapan/internal/adapter/rest/echo/handler/model"
	"github.com/ramadani/balapan/internal/app/command"
	rmodel "github.com/ramadani/balapan/internal/app/command/model"
	"net/http"
)

type RewardsHandler struct {
	usageCommand command.UsageRewardsCommander
}

func (h *RewardsHandler) Usage(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	req := new(model.UsageRewardsRequest)

	if err := c.Bind(&req); err != nil {
		return err
	}

	data := &rmodel.UsageRewards{
		ID:     id,
		Amount: req.Amount,
	}

	if err := h.usageCommand.Do(ctx, data); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func NewRewardsHandler(usageCommand command.UsageRewardsCommander) *RewardsHandler {
	return &RewardsHandler{usageCommand: usageCommand}
}
