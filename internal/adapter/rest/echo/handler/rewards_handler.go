package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/ramadani/balapan/internal/adapter/rest/echo/handler/model"
	"github.com/ramadani/balapan/internal/app/command"
	rmodel "github.com/ramadani/balapan/internal/app/command/model"
	"net/http"
)

type RewardsHandler struct {
	claimCommand command.ClaimRewardsCommander
}

func (h *RewardsHandler) Claim(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	req := new(model.ClaimRewardsRequest)

	if err := c.Bind(&req); err != nil {
		return err
	}

	data := &rmodel.ClaimRewards{
		ID:     id,
		UserID: req.UserID,
		Amount: req.Amount,
	}

	if err := h.claimCommand.Do(ctx, data); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func NewRewardsHandler(claimCommand command.ClaimRewardsCommander) *RewardsHandler {
	return &RewardsHandler{claimCommand: claimCommand}
}
