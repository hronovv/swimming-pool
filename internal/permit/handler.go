package permit

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"pool/internal/auth"
	"pool/internal/user"
)

type PermitService interface {
	GetAll(ctx context.Context, filter Filter) ([]Permit, error)
	Create(ctx context.Context, item Permit) (Permit, error)
	TrainerIDByUser(ctx context.Context, userID int) (int, error)
}

type Handler struct {
	service PermitService
}

func NewHandler(service PermitService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/permits/", h.GetAll)
	r.POST("/permits/", h.Create)
}

type PermitRequest struct {
	TrainerID    int    `json:"trainer_id"`
	WorkFormatID int    `json:"work_format_id"`
	IssuedOn     string `json:"issued_on"`
	ExpiresOn    string `json:"expires_on"`
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
	if actor, found := auth.CurrentUser(c); found && actor.Role == user.RoleTrainer {
		own, err := h.service.TrainerIDByUser(c.Request.Context(), actor.ID)
		if err != nil {
			writePermitError(c, err)
			return
		}
		if own == 0 {
			c.JSON(http.StatusOK, []PermitResponse{})
			return
		}
		if trainerID != 0 && trainerID != own {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		trainerID = own
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

func (h *Handler) Create(c *gin.Context) {
	var req PermitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	issuedOn, ok := parseDate(c, "issued_on", req.IssuedOn, true)
	if !ok {
		return
	}
	var expiresOn *time.Time
	if strings.TrimSpace(req.ExpiresOn) != "" {
		expires, valid := parseDate(c, "expires_on", req.ExpiresOn, true)
		if !valid {
			return
		}
		expiresOn = &expires
	}
	created, err := h.service.Create(c.Request.Context(), Permit{
		TrainerID:    req.TrainerID,
		WorkFormatID: req.WorkFormatID,
		IssuedOn:     issuedOn,
		ExpiresOn:    expiresOn,
	})
	if err != nil {
		writePermitError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toPermitResponse(created))
}

func parseDate(c *gin.Context, field, raw string, required bool) (time.Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		if required {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permit: " + field + " is required"})
		}
		return time.Time{}, !required
	}
	parsed, err := time.Parse("2006-01-02", raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permit: " + field + " must be YYYY-MM-DD"})
		return time.Time{}, false
	}
	return parsed, true
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
	case errors.Is(err, ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": "permit already exists"})
	case errors.Is(err, ErrTrainerNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "trainer not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
