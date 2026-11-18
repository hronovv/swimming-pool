package client

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientService interface {
	GetAll(ctx context.Context, filter Filter) ([]Client, error)
	GetByID(ctx context.Context, ID int) (Client, error)
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
	}
}

type ClientResponse struct {
	ID                 int     `json:"id"`
	UserID             int     `json:"user_id"`
	Passport           string  `json:"passport"`
	BirthDate          string  `json:"birth_date"`
	Address            *string `json:"address"`
	MedicalCertificate bool    `json:"medical_certificate"`
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
		resp = append(resp, toClientResponse(item))
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
	c.JSON(http.StatusOK, toClientResponse(item))
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

func toClientResponse(item Client) ClientResponse {
	return ClientResponse{
		ID:                 item.ID,
		UserID:             item.UserID,
		Passport:           item.Passport,
		BirthDate:          item.BirthDate.Format("02.01.2006"),
		Address:            item.Address,
		MedicalCertificate: item.MedicalCertificate,
	}
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
