package tariff

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TariffService interface {
	GetAll(ctx context.Context) ([]Tariff, error)
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
	users := r.Group("/tariffs")
	{
		users.GET("/", h.GetAll)
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	tariffs, err := h.service.GetAll(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, tariffs)

}
