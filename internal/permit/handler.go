package permit

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PermitService interface {
	GetAll(ctx context.Context, filter Filter) ([]Permit, error)
}

type Handler struct {
	service PermitService
}

func NewHandler(service PermitService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/permits/", h.GetAll)
}

type PermitResponse struct {
	TrainerID    int     `json:"trainer_id"`
	WorkFormatID int     `json:"work_format_id"`
	IssuedOn     string  `json:"issued_on"`
	ExpiresOn    *string `json:"expires_on"`
}

func (h *Handler) GetAll(c *gin.Context) {
	trainerID, ok := queryID(c, "trainer_id")
	if !ok {
		return
	}
	items, err := h.service.GetAll(c.Request.Context(), Filter{TrainerID: trainerID})
	if err != nil {
		writePermitError(c, err)
		return
	}
	resp := make([]PermitResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, toPermitResponse(item))
	}
	c.JSON(http.StatusOK, resp)
}

func queryID(c *gin.Context, name string) (int, bool) {
	raw := c.Query(name)
	if raw == "" {
		return 0, true
	}
	id, err := strconv.Atoi(raw)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func toPermitResponse(item Permit) PermitResponse {
	resp := PermitResponse{
		TrainerID:    item.TrainerID,
		WorkFormatID: item.WorkFormatID,
		IssuedOn:     item.IssuedOn.Format("02.01.2006"),
	}
	if item.ExpiresOn != nil {
		value := item.ExpiresOn.Format("02.01.2006")
		resp.ExpiresOn = &value
	}
	return resp
}

func writePermitError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
