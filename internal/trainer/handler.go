package trainer

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TrainerService interface {
	GetAll(ctx context.Context, filter Filter) ([]Trainer, error)
	GetByID(ctx context.Context, ID int) (Trainer, error)
}

type Handler struct {
	service TrainerService
}

func NewHandler(service TrainerService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	trainers := r.Group("/trainers")
	{
		trainers.GET("/", h.GetAll)
		trainers.GET("/:id", h.GetByID)
	}
}

type TrainerResponse struct {
	ID              int    `json:"id"`
	UserID          int    `json:"user_id"`
	Qualification   string `json:"qualification"`
	Rank            string `json:"rank"`
	ExperienceYears int    `json:"experience_years"`
	Specialization  string `json:"specialization"`
}

func (h *Handler) GetAll(c *gin.Context) {
	userID, ok := queryID(c, "user_id")
	if !ok {
		return
	}
	items, err := h.service.GetAll(c.Request.Context(), Filter{UserID: userID})
	if err != nil {
		writeTrainerError(c, err)
		return
	}
	resp := make([]TrainerResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, toTrainerResponse(item))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		writeTrainerError(c, err)
		return
	}
	c.JSON(http.StatusOK, toTrainerResponse(item))
}

func parseID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return 0, false
	}
	return id, true
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

func toTrainerResponse(item Trainer) TrainerResponse {
	return TrainerResponse{
		ID:              item.ID,
		UserID:          item.UserID,
		Qualification:   item.Qualification,
		Rank:            item.Rank,
		ExperienceYears: item.ExperienceYears,
		Specialization:  item.Specialization,
	}
}

func writeTrainerError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "trainer not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
