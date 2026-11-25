package group

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type GroupService interface {
	GetAll(ctx context.Context, filter Filter) ([]Group, error)
	GetByID(ctx context.Context, ID int) (Group, error)
	Members(ctx context.Context, groupID int) ([]Member, error)
	Create(ctx context.Context, group Group) (Group, error)
	Update(ctx context.Context, group Group) error
	UpdateStatus(ctx context.Context, id int, status string) error
}

type Handler struct {
	service GroupService
}

func NewHandler(service GroupService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	groups := r.Group("/groups")
	{
		groups.GET("/", h.GetAll)
		groups.POST("/", h.Create)
		groups.GET("/:id/members", h.Members)
		groups.GET("/:id", h.GetByID)
		groups.PUT("/:id", h.Update)
		groups.PATCH("/:id/status", h.UpdateStatus)
	}
}

type MemberResponse struct {
	ClientID      int    `json:"client_id"`
	Login         string `json:"login"`
	Name          string `json:"name"`
	ServiceID     int    `json:"service_id"`
	ServiceStatus string `json:"service_status"`
}

type GroupResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Capacity  int    `json:"capacity"`
	Status    string `json:"status"`
}

type GroupRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Capacity  int    `json:"capacity"`
	Status    string `json:"status"`
}

type StatusRequest struct {
	Status string `json:"status"`
}

func (h *Handler) GetAll(c *gin.Context) {
	groups, err := h.service.GetAll(c.Request.Context(), Filter{Status: c.Query("status")})
	if err != nil {
		writeGroupError(c, err)
		return
	}

	resp := make([]GroupResponse, 0, len(groups))
	for _, group := range groups {
		resp = append(resp, toGroupResponse(group))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	group, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		writeGroupError(c, err)
		return
	}

	c.JSON(http.StatusOK, toGroupResponse(group))
}

func (h *Handler) Members(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	members, err := h.service.Members(c.Request.Context(), id)
	if err != nil {
		writeGroupError(c, err)
		return
	}

	resp := make([]MemberResponse, 0, len(members))
	for _, member := range members {
		resp = append(resp, toMemberResponse(member))
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Create(c *gin.Context) {
	var req GroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	group, ok := toGroup(c, req, 0)
	if !ok {
		return
	}

	created, err := h.service.Create(c.Request.Context(), group)
	if err != nil {
		writeGroupError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toGroupResponse(created))
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req GroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	group, ok := toGroup(c, req, id)
	if !ok {
		return
	}

	if err := h.service.Update(c.Request.Context(), group); err != nil {
		writeGroupError(c, err)
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
		writeGroupError(c, err)
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

func toGroup(c *gin.Context, req GroupRequest, id int) (Group, bool) {
	startDate, ok := parseDate(c, req.StartDate, "start_date")
	if !ok {
		return Group{}, false
	}
	endDate, ok := parseDate(c, req.EndDate, "end_date")
	if !ok {
		return Group{}, false
	}
	return Group{
		ID:        id,
		Name:      req.Name,
		StartDate: startDate,
		EndDate:   endDate,
		Capacity:  req.Capacity,
		Status:    req.Status,
	}, true
}

func toGroupResponse(group Group) GroupResponse {
	return GroupResponse{
		ID:        group.ID,
		Name:      group.Name,
		StartDate: group.StartDate.Format("02.01.2006"),
		EndDate:   group.EndDate.Format("02.01.2006"),
		Capacity:  group.Capacity,
		Status:    group.Status,
	}
}

func toMemberResponse(member Member) MemberResponse {
	name := member.LastName + " " + firstInitial(member.FirstName) + "."
	if member.Patronymic != nil && *member.Patronymic != "" {
		name += " " + firstInitial(*member.Patronymic) + "."
	}
	return MemberResponse{
		ClientID:      member.ClientID,
		Login:         member.Login,
		Name:          name,
		ServiceID:     member.ServiceID,
		ServiceStatus: member.ServiceStatus,
	}
}

func firstInitial(s string) string {
	r := []rune(s)
	if len(r) == 0 {
		return ""
	}
	return string(r[0])
}

func writeGroupError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
