package service

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"pool/internal/auth"
	"pool/internal/user"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	services := r.Group("/services")
	{
		services.GET("/", h.GetAll)
		services.POST("/", h.Create)
		services.GET("/:id", h.GetByID)
		services.PUT("/:id", h.Update)
		services.PATCH("/:id/status", h.UpdateStatus)
		services.DELETE("/:id", h.Delete)
	}
}

type ServiceResponse struct {
	ID           int     `json:"id"`
	ClientID     int     `json:"client_id"`
	TariffID     int     `json:"tariff_id"`
	GroupID      *int    `json:"group_id"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	Status       string  `json:"status"`
	ClientName   string  `json:"client_name"`
	TariffName   string  `json:"tariff_name"`
	WorkFormatID int     `json:"work_format_id"`
	GroupName    *string `json:"group_name"`
}

type ServiceRequest struct {
	ClientID  int    `json:"client_id"`
	TariffID  int    `json:"tariff_id"`
	GroupID   *int   `json:"group_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Status    string `json:"status"`
}

type StatusRequest struct {
	Status string `json:"status"`
}

func (h *Handler) GetAll(c *gin.Context) {
	filter, ok := parseFilter(c)
	if !ok {
		return
	}
	filter, ok = h.ownServices(c, filter)
	if !ok {
		return
	}

	models, err := h.service.GetAll(c.Request.Context(), filter)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	resp := make([]ServiceResponse, 0, len(models))
	for _, model := range models {
		resp = append(resp, toServiceResponse(model))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	model, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	if !h.canSeeService(c, model.ClientID) {
		return
	}

	c.JSON(http.StatusOK, toServiceResponse(model))
}

func (h *Handler) ownServices(c *gin.Context, filter Filter) (Filter, bool) {
	actor, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return filter, false
	}
	if actor.Role == user.RoleTrainer {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return filter, false
	}
	if actor.Role != user.RoleClient {
		return filter, true
	}
	clientID, err := h.service.ClientIDByUser(c.Request.Context(), actor.ID)
	if err != nil {
		writeServiceError(c, err)
		return filter, false
	}
	if clientID == 0 {
		c.JSON(http.StatusOK, []ServiceResponse{})
		return filter, false
	}
	filter.ClientID = clientID
	return filter, true
}

func (h *Handler) canSeeService(c *gin.Context, clientID int) bool {
	actor, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false
	}
	if actor.Role == user.RoleTrainer {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return false
	}
	if actor.Role != user.RoleClient {
		return true
	}
	own, err := h.service.ClientIDByUser(c.Request.Context(), actor.ID)
	if err != nil {
		writeServiceError(c, err)
		return false
	}
	if own != clientID {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return false
	}
	return true
}

func (h *Handler) Create(c *gin.Context) {
	var req ServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	model, ok := toModel(c, req, 0)
	if !ok {
		return
	}

	created, err := h.service.Create(c.Request.Context(), model)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toServiceResponse(created))
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req ServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	model, ok := toModel(c, req, id)
	if !ok {
		return
	}

	if err := h.service.Update(c.Request.Context(), model); err != nil {
		writeServiceError(c, err)
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
		writeServiceError(c, err)
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
		writeServiceError(c, err)
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

func parseFilter(c *gin.Context) (Filter, bool) {
	filter := Filter{Status: c.Query("status")}
	var ok bool
	if filter.ClientID, ok = queryID(c, "client_id"); !ok {
		return Filter{}, false
	}
	if filter.TariffID, ok = queryID(c, "tariff_id"); !ok {
		return Filter{}, false
	}
	if filter.GroupID, ok = queryID(c, "group_id"); !ok {
		return Filter{}, false
	}
	return filter, true
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
		return time.Time{}, true
	}
	value, err := time.Parse("2006-01-02", raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": field + " must be YYYY-MM-DD"})
		return time.Time{}, false
	}
	return value, true
}

func toModel(c *gin.Context, req ServiceRequest, id int) (Model, bool) {
	startDate, ok := parseDate(c, req.StartDate, "start_date")
	if !ok {
		return Model{}, false
	}
	endDate, ok := parseDate(c, req.EndDate, "end_date")
	if !ok {
		return Model{}, false
	}
	return Model{
		ID:        id,
		ClientID:  req.ClientID,
		TariffID:  req.TariffID,
		GroupID:   req.GroupID,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    req.Status,
	}, true
}

func toServiceResponse(model Model) ServiceResponse {
	return ServiceResponse{
		ID:           model.ID,
		ClientID:     model.ClientID,
		TariffID:     model.TariffID,
		GroupID:      model.GroupID,
		StartDate:    model.StartDate.Format("02.01.2006"),
		EndDate:      model.EndDate.Format("02.01.2006"),
		Status:       model.Status,
		ClientName:   model.ClientName,
		TariffName:   model.TariffName,
		WorkFormatID: model.WorkFormatID,
		GroupName:    model.GroupName,
	}
}

func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
	case errors.Is(err, ErrClientNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
	case errors.Is(err, ErrTariffNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "tariff not found"})
	case errors.Is(err, ErrGroupNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
	case errors.Is(err, ErrInUse):
		c.JSON(http.StatusConflict, gin.H{"error": "service is still in use"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
