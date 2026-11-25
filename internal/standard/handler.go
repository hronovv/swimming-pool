package standard

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"pool/internal/auth"
	"pool/internal/user"
)

type StandardService interface {
	GetAll(ctx context.Context, filter Filter) ([]Standard, error)
	GetByID(ctx context.Context, ID int) (Standard, error)
	ClientIDByUser(ctx context.Context, userID int) (int, error)
	Create(ctx context.Context, item Standard) (Standard, error)
	Update(ctx context.Context, item Standard) error
	Delete(ctx context.Context, id int) error
	TrainerIDByUser(ctx context.Context, userID int) (int, error)
}

type Handler struct {
	service StandardService
}

func NewHandler(service StandardService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	standards := r.Group("/standards")
	{
		standards.GET("/", h.GetAll)
		standards.POST("/", h.Create)
		standards.GET("/:id", h.GetByID)
		standards.PUT("/:id", h.Update)
		standards.DELETE("/:id", h.Delete)
	}
}

type StandardResponse struct {
	ID          int     `json:"id"`
	ClientID    int     `json:"client_id"`
	TrainerID   int     `json:"trainer_id"`
	Style       string  `json:"style"`
	DistanceM   int     `json:"distance_m"`
	HeldOn      string  `json:"held_on"`
	ResultTime  *string `json:"result_time"`
	Outcome     string  `json:"outcome"`
	AttemptNo   int     `json:"attempt_no"`
	Comment     *string `json:"comment"`
	ClientName  string  `json:"client_name"`
	TrainerName string  `json:"trainer_name"`
}

type StandardRequest struct {
	ClientID   int    `json:"client_id"`
	TrainerID  int    `json:"trainer_id"`
	Style      string `json:"style"`
	DistanceM  int    `json:"distance_m"`
	HeldOn     string `json:"held_on"`
	ResultTime string `json:"result_time"`
	Outcome    string `json:"outcome"`
	AttemptNo  int    `json:"attempt_no"`
	Comment    string `json:"comment"`
}

func (h *Handler) GetAll(c *gin.Context) {
	filter := Filter{}
	var ok bool
	if filter.ClientID, ok = queryID(c, "client_id"); !ok {
		return
	}
	if filter.TrainerID, ok = queryID(c, "trainer_id"); !ok {
		return
	}
	if actor, found := auth.CurrentUser(c); found {
		switch actor.Role {
		case user.RoleClient:
			clientID, err := h.service.ClientIDByUser(c.Request.Context(), actor.ID)
			if err != nil {
				writeStandardError(c, err)
				return
			}
			if clientID == 0 {
				c.JSON(http.StatusOK, []StandardResponse{})
				return
			}
			filter.ClientID = clientID
		case user.RoleTrainer:
			trainerID, err := h.service.TrainerIDByUser(c.Request.Context(), actor.ID)
			if err != nil {
				writeStandardError(c, err)
				return
			}
			if trainerID == 0 {
				c.JSON(http.StatusOK, []StandardResponse{})
				return
			}
			if filter.TrainerID != 0 && filter.TrainerID != trainerID {
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
			filter.TrainerID = trainerID
		}
	}
	items, err := h.service.GetAll(c.Request.Context(), filter)
	if err != nil {
		writeStandardError(c, err)
		return
	}
	resp := make([]StandardResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, toStandardResponse(item))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	item, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		writeStandardError(c, err)
		return
	}
	if actor, found := auth.CurrentUser(c); found {
		switch actor.Role {
		case user.RoleClient:
			clientID, err := h.service.ClientIDByUser(c.Request.Context(), actor.ID)
			if err != nil {
				writeStandardError(c, err)
				return
			}
			if clientID != item.ClientID {
				c.JSON(http.StatusNotFound, gin.H{"error": "standard not found"})
				return
			}
		case user.RoleTrainer:
			trainerID, err := h.service.TrainerIDByUser(c.Request.Context(), actor.ID)
			if err != nil {
				writeStandardError(c, err)
				return
			}
			if trainerID != item.TrainerID {
				c.JSON(http.StatusNotFound, gin.H{"error": "standard not found"})
				return
			}
		}
	}
	c.JSON(http.StatusOK, toStandardResponse(item))
}

func (h *Handler) Create(c *gin.Context) {
	item, ok := h.parseStandard(c, 0)
	if !ok {
		return
	}
	created, err := h.service.Create(c.Request.Context(), item)
	if err != nil {
		writeStandardError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toStandardResponse(created))
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	item, ok := h.parseStandard(c, id)
	if !ok {
		return
	}
	if err := h.service.Update(c.Request.Context(), item); err != nil {
		writeStandardError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		writeStandardError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) parseStandard(c *gin.Context, id int) (Standard, bool) {
	var req StandardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return Standard{}, false
	}
	heldOn, ok := parseDate(c, req.HeldOn, "held_on")
	if !ok {
		return Standard{}, false
	}
	result, ok := parseClock(c, req.ResultTime)
	if !ok {
		return Standard{}, false
	}
	var comment *string
	if req.Comment != "" {
		comment = &req.Comment
	}
	trainerID := req.TrainerID
	if actor, found := auth.CurrentUser(c); found && actor.Role == user.RoleTrainer {
		own, err := h.service.TrainerIDByUser(c.Request.Context(), actor.ID)
		if err != nil {
			writeStandardError(c, err)
			return Standard{}, false
		}
		if own == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return Standard{}, false
		}
		trainerID = own
	}
	return Standard{
		ID:         id,
		ClientID:   req.ClientID,
		TrainerID:  trainerID,
		Style:      req.Style,
		DistanceM:  req.DistanceM,
		HeldOn:     heldOn,
		ResultTime: result,
		Outcome:    req.Outcome,
		AttemptNo:  req.AttemptNo,
		Comment:    comment,
	}, true
}

func parseID(c *gin.Context, name string) (int, bool) {
	id, err := strconv.Atoi(c.Param(name))
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

func parseDate(c *gin.Context, raw, field string) (time.Time, bool) {
	if raw == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": field + " is required"})
		return time.Time{}, false
	}
	value, err := time.Parse("2006-01-02", raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": field + " must be YYYY-MM-DD"})
		return time.Time{}, false
	}
	return value, true
}

func parseClock(c *gin.Context, raw string) (*time.Time, bool) {
	if raw == "" {
		return nil, true
	}
	for _, layout := range []string{"15:04:05", "15:04"} {
		value, err := time.Parse(layout, raw)
		if err == nil {
			return &value, true
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "result_time must be HH:MM:SS"})
	return nil, false
}

func toStandardResponse(item Standard) StandardResponse {
	resp := StandardResponse{
		ID:          item.ID,
		ClientID:    item.ClientID,
		TrainerID:   item.TrainerID,
		Style:       item.Style,
		DistanceM:   item.DistanceM,
		HeldOn:      item.HeldOn.Format("02.01.2006"),
		Outcome:     item.Outcome,
		AttemptNo:   item.AttemptNo,
		Comment:     item.Comment,
		ClientName:  item.ClientName,
		TrainerName: item.TrainerName,
	}
	if item.ResultTime != nil {
		value := item.ResultTime.Format("15:04:05")
		resp.ResultTime = &value
	}
	return resp
}

func writeStandardError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "standard not found"})
	case errors.Is(err, ErrClientNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
	case errors.Is(err, ErrTrainerNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "trainer not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
