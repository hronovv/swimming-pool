package enrollment

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pool/internal/auth"
	"pool/internal/user"
)

type EnrollmentService interface {
	GetAll(ctx context.Context, filter Filter) ([]Enrollment, error)
	Create(ctx context.Context, item Enrollment) (Enrollment, error)
	UpdateAttended(ctx context.Context, serviceID, sessionID int, attended bool) error
	Delete(ctx context.Context, serviceID, sessionID int) error
}

type Handler struct {
	service EnrollmentService
}

func NewHandler(service EnrollmentService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	enrollments := r.Group("/enrollments")
	{
		enrollments.GET("/", h.GetAll)
		enrollments.POST("/", h.Create)
		enrollments.PATCH("/:service_id/:session_id", h.UpdateAttended)
		enrollments.DELETE("/:service_id/:session_id", h.Delete)
	}
}

type EnrollmentResponse struct {
	ServiceID  int     `json:"service_id"`
	SessionID  int     `json:"session_id"`
	Attended   bool    `json:"attended"`
	Comment    *string `json:"comment"`
	ClientID   int     `json:"client_id"`
	ClientName string  `json:"client_name"`
	TariffName string  `json:"tariff_name"`
}

type EnrollmentRequest struct {
	ServiceID int    `json:"service_id"`
	SessionID int    `json:"session_id"`
	Attended  bool   `json:"attended"`
	Comment   string `json:"comment"`
}

type AttendedRequest struct {
	Attended *bool `json:"attended"`
}

func (h *Handler) GetAll(c *gin.Context) {
	filter := Filter{}
	var ok bool
	if filter.ServiceID, ok = queryID(c, "service_id"); !ok {
		return
	}
	if filter.SessionID, ok = queryID(c, "session_id"); !ok {
		return
	}
	if actor, ok := auth.CurrentUser(c); ok && actor.Role == user.RoleClient {
		filter.ClientUserID = actor.ID
	}
	items, err := h.service.GetAll(c.Request.Context(), filter)
	if err != nil {
		writeEnrollmentError(c, err)
		return
	}
	resp := make([]EnrollmentResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, toEnrollmentResponse(item))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Create(c *gin.Context) {
	var req EnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	var comment *string
	if req.Comment != "" {
		comment = &req.Comment
	}
	item := Enrollment{
		ServiceID: req.ServiceID,
		SessionID: req.SessionID,
		Attended:  req.Attended,
		Comment:   comment,
	}
	if actor, ok := auth.CurrentUser(c); ok && actor.Role == user.RoleClient {
		item.ActorUserID = actor.ID
	}
	created, err := h.service.Create(c.Request.Context(), item)
	if err != nil {
		writeEnrollmentError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toEnrollmentResponse(created))
}

func (h *Handler) UpdateAttended(c *gin.Context) {
	serviceID, ok := pathID(c, "service_id")
	if !ok {
		return
	}
	sessionID, ok := pathID(c, "session_id")
	if !ok {
		return
	}
	var req AttendedRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Attended == nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if err := h.service.UpdateAttended(c.Request.Context(), serviceID, sessionID, *req.Attended); err != nil {
		writeEnrollmentError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Delete(c *gin.Context) {
	serviceID, ok := pathID(c, "service_id")
	if !ok {
		return
	}
	sessionID, ok := pathID(c, "session_id")
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), serviceID, sessionID); err != nil {
		writeEnrollmentError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func pathID(c *gin.Context, name string) (int, bool) {
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

func toEnrollmentResponse(item Enrollment) EnrollmentResponse {
	return EnrollmentResponse{
		ServiceID:  item.ServiceID,
		SessionID:  item.SessionID,
		Attended:   item.Attended,
		Comment:    item.Comment,
		ClientID:   item.ClientID,
		ClientName: item.ClientName,
		TariffName: item.TariffName,
	}
}

func writeEnrollmentError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "enrollment not found"})
	case errors.Is(err, ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": "enrollment already exists"})
	case errors.Is(err, ErrServiceNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
	case errors.Is(err, ErrSessionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
