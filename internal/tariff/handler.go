package tariff

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TariffService interface {
	GetAll(ctx context.Context) ([]Tariff, error)
	GetByID(ctx context.Context, ID int) (Tariff, error)
	Create(ctx context.Context, tariff Tariff) (Tariff, error)
	Update(ctx context.Context, tariff Tariff) error
	Delete(ctx context.Context, id int) error
}

type Handler struct {
	service TariffService
}

func NewHandler(service TariffService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	tariffs := r.Group("/tariffs")
	{
		tariffs.GET("/", h.GetAll)
		tariffs.POST("/", h.Create)
		tariffs.GET("/:id", h.GetByID)
		tariffs.PUT("/:id", h.Update)
		tariffs.DELETE("/:id", h.Delete)
	}
}

type TariffResponse struct {
	ID                   int     `json:"id"`
	Name                 string  `json:"name"`
	Format               string  `json:"format"`
	WorkFormatID         int     `json:"work_format_id"`
	Price                float64 `json:"price"`
	VisitLimit           *int    `json:"visit_limit"`
	DurationDays         int     `json:"duration_days"`
	ServicesCnt          int64   `json:"services_cnt"`
	CompletedPaymentsCnt int64   `json:"completed_payments_cnt"`
	Revenue              float64 `json:"revenue"`
	AvgCheck             float64 `json:"avg_check"`
}

type TariffRequest struct {
	Name         string  `json:"name"`
	Format       string  `json:"format"`
	Price        float64 `json:"price"`
	VisitLimit   *int    `json:"visit_limit"`
	DurationDays int     `json:"duration_days"`
}

func (h *Handler) GetAll(c *gin.Context) {
	tariffs, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		writeTariffError(c, err)
		return
	}

	resp := make([]TariffResponse, 0, len(tariffs))
	for _, tariff := range tariffs {
		resp = append(resp, toTariffResponse(tariff))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	tariff, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		writeTariffError(c, err)
		return
	}

	c.JSON(http.StatusOK, toTariffResponse(tariff))
}

func (h *Handler) Create(c *gin.Context) {
	var req TariffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	created, err := h.service.Create(c.Request.Context(), toTariff(req, 0))
	if err != nil {
		writeTariffError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toTariffResponse(created))
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req TariffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.service.Update(c.Request.Context(), toTariff(req, id)); err != nil {
		writeTariffError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		writeTariffError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func toTariff(req TariffRequest, id int) Tariff {
	return Tariff{
		ID:           id,
		Name:         req.Name,
		Format:       req.Format,
		Price:        req.Price,
		VisitLimit:   req.VisitLimit,
		DurationDays: req.DurationDays,
	}
}

func toTariffResponse(tariff Tariff) TariffResponse {
	return TariffResponse{
		ID:                   tariff.ID,
		Name:                 tariff.Name,
		Format:               tariff.Format,
		WorkFormatID:         tariff.WorkFormatID,
		Price:                tariff.Price,
		VisitLimit:           tariff.VisitLimit,
		DurationDays:         tariff.DurationDays,
		ServicesCnt:          tariff.ServicesCnt,
		CompletedPaymentsCnt: tariff.CompletedPaymentsCnt,
		Revenue:              tariff.Revenue,
		AvgCheck:             tariff.AvgCheck,
	}
}

func writeTariffError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "tariff not found"})
	case errors.Is(err, ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": "tariff already exists"})
	case errors.Is(err, ErrInUse):
		c.JSON(http.StatusConflict, gin.H{"error": "tariff is still in use"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
