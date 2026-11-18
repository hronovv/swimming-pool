package lane

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LaneService interface {
	GetAll(ctx context.Context) ([]Lane, error)
	GetByID(ctx context.Context, ID int) (Lane, error)
	Create(ctx context.Context, lane Lane) (Lane, error)
	Update(ctx context.Context, lane Lane) error
	UpdateStatus(ctx context.Context, id int, status string) error
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
		lanes.POST("/", h.Create)
		lanes.GET("/:id", h.GetByID)
		lanes.PUT("/:id", h.Update)
		lanes.PATCH("/:id/status", h.UpdateStatus)
	}
}

type LaneResponse struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	LengthM int     `json:"length_m"`
	DepthM  float64 `json:"depth_m"`
	Status  string  `json:"status"`
}

type LaneRequest struct {
	Name    string  `json:"name"`
	LengthM int     `json:"length_m"`
	DepthM  float64 `json:"depth_m"`
	Status  string  `json:"status"`
}

type StatusRequest struct {
	Status string `json:"status"`
}

func (h *Handler) GetAll(c *gin.Context) {
	lanes, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		writeLaneError(c, err)
		return
	}

	resp := make([]LaneResponse, 0, len(lanes))
	for _, lane := range lanes {
		resp = append(resp, toLaneResponse(lane))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	lane, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		writeLaneError(c, err)
		return
	}

	c.JSON(http.StatusOK, toLaneResponse(lane))
}

func (h *Handler) Create(c *gin.Context) {
	var req LaneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	created, err := h.service.Create(c.Request.Context(), Lane{
		Name:    req.Name,
		LengthM: req.LengthM,
		DepthM:  req.DepthM,
		Status:  req.Status,
	})
	if err != nil {
		writeLaneError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toLaneResponse(created))
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req LaneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err := h.service.Update(c.Request.Context(), Lane{
		ID:      id,
		Name:    req.Name,
		LengthM: req.LengthM,
		DepthM:  req.DepthM,
		Status:  req.Status,
	})
	if err != nil {
		writeLaneError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		writeLaneError(c, err)
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

func toLaneResponse(lane Lane) LaneResponse {
	return LaneResponse{
		ID:      lane.ID,
		Name:    lane.Name,
		LengthM: lane.LengthM,
		DepthM:  lane.DepthM,
		Status:  lane.Status,
	}
}

func writeLaneError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "lane not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
