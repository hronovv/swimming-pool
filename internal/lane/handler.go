package lane

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LaneService interface {
	GetAll(ctx context.Context) ([]Lane, error)
	GetByID(ctx context.Context, ID int) (Lane, error)
}

type Handler struct {
	service LaneService
}

func NewHandler(service LaneService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	lanes := r.Group("/lanes")
	{
		lanes.GET("/", h.GetAll)
		lanes.GET("/:id", h.GetByID)

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

type LaneResponse struct {
	Name    string  `json:"name"`
	LengthM int     `json:"length_m"`
	DepthM  float64 `json:"depth_m"`
	Status  string  `json:"status"`
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	lane, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	resp := LaneResponse{
		Name:    lane.Name,
		LengthM: lane.LengthM,
		DepthM:  lane.DepthM,
		Status:  lane.Status,
	}

	c.JSON(http.StatusOK, resp)

}
