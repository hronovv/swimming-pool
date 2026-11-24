package client

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pool/internal/auth"
	"pool/internal/user"
)

type ClientService interface {
	GetAll(ctx context.Context, filter Filter) ([]Client, error)
	GetByID(ctx context.Context, ID int) (Client, error)
	UpdateMedical(ctx context.Context, id int, allowed bool) error
}

type Handler struct {
	service ClientService
}

func NewHandler(service ClientService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	clients := r.Group("/clients")
	{
		clients.GET("/", h.GetAll)
		clients.GET("/:id", h.GetByID)
		clients.PATCH("/:id/medical", h.UpdateMedical)
	}
}

type ClientResponse struct {
	ID                 int     `json:"id"`
	UserID             int     `json:"user_id"`
	Passport           string  `json:"passport"`
	BirthDate          string  `json:"birth_date"`
	Address            *string `json:"address"`
	MedicalCertificate bool    `json:"medical_certificate"`
	Name               string  `json:"name"`
}

type MedicalRequest struct {
	MedicalCertificate *bool `json:"medical_certificate"`
}

func (h *Handler) GetAll(c *gin.Context) {
	userID, ok := queryID(c, "user_id")
	if !ok {
		return
	}
	items, err := h.service.GetAll(c.Request.Context(), Filter{UserID: userID})
	if err != nil {
		writeClientError(c, err)
		return
	}
	resp := make([]ClientResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, h.visible(c, item))
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
		writeClientError(c, err)
		return
	}
	c.JSON(http.StatusOK, h.visible(c, item))
}

func (h *Handler) UpdateMedical(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req MedicalRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.MedicalCertificate == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "medical_certificate is required"})
		return
	}
	if err := h.service.UpdateMedical(c.Request.Context(), id, *req.MedicalCertificate); err != nil {
		writeClientError(c, err)
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

func (h *Handler) visible(c *gin.Context, item Client) ClientResponse {
	resp := ClientResponse{
		ID:                 item.ID,
		UserID:             item.UserID,
		Passport:           item.Passport,
		BirthDate:          item.BirthDate.Format("02.01.2006"),
		Address:            item.Address,
		MedicalCertificate: item.MedicalCertificate,
		Name:               item.Name,
	}
	actor, ok := auth.CurrentUser(c)
	if ok && actor.Role != user.RoleAdmin {
		resp.Passport = ""
		resp.BirthDate = ""
		resp.Address = nil
		resp.MedicalCertificate = false
	}
	return resp
}

func writeClientError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
