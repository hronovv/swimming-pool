package session

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SessionService interface {
	GetAll(ctx context.Context, filter Filter) ([]Session, error)
	GetByID(ctx context.Context, ID int) (Session, error)
	Create(ctx context.Context, session Session) (Session, error)
	Update(ctx context.Context, session Session) error
	UpdateHeld(ctx context.Context, id int, wasHeld bool) error
}

type Handler struct {
	service SessionService
}

func NewHandler(service SessionService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	sessions := r.Group("/sessions")
	{
		sessions.GET("/", h.GetAll)
		sessions.POST("/", h.Create)
		sessions.GET("/:id", h.GetByID)
		sessions.PUT("/:id", h.Update)
		sessions.PATCH("/:id/held", h.UpdateHeld)
	}
}

type SessionResponse struct {
	ID           int     `json:"id"`
	WorkFormatID int     `json:"work_format_id"`
	GroupID      *int    `json:"group_id"`
	LaneID       *int    `json:"lane_id"`
	SessionDate  string  `json:"session_date"`
	StartTime    string  `json:"start_time"`
	EndTime      string  `json:"end_time"`
	Topic        *string `json:"topic"`
	WasHeld      bool    `json:"was_held"`
	TrainerIDs   []int   `json:"trainer_ids"`
}

type SessionRequest struct {
	WorkFormatID int     `json:"work_format_id"`
	GroupID      *int    `json:"group_id"`
	LaneID       *int    `json:"lane_id"`
	SessionDate  string  `json:"session_date"`
	StartTime    string  `json:"start_time"`
	EndTime      string  `json:"end_time"`
	Topic        *string `json:"topic"`
	WasHeld      bool    `json:"was_held"`
	TrainerIDs   []int   `json:"trainer_ids"`
}

type HeldRequest struct {
	WasHeld *bool `json:"was_held"`
}

func (h *Handler) GetAll(c *gin.Context) {
	filter, ok := parseFilter(c)
	if !ok {
		return
	}

	sessions, err := h.service.GetAll(c.Request.Context(), filter)
	if err != nil {
		writeSessionError(c, err)
		return
	}

	resp := make([]SessionResponse, 0, len(sessions))
	for _, item := range sessions {
		resp = append(resp, toSessionResponse(item))
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
		writeSessionError(c, err)
		return
	}

	c.JSON(http.StatusOK, toSessionResponse(item))
}

func (h *Handler) Create(c *gin.Context) {
	var req SessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	item, ok := toSession(c, req, 0)
	if !ok {
		return
	}

	created, err := h.service.Create(c.Request.Context(), item)
	if err != nil {
		writeSessionError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toSessionResponse(created))
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req SessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	item, ok := toSession(c, req, id)
	if !ok {
		return
	}

	if err := h.service.Update(c.Request.Context(), item); err != nil {
		writeSessionError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdateHeld(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req HeldRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.WasHeld == nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateHeld(c.Request.Context(), id, *req.WasHeld); err != nil {
		writeSessionError(c, err)
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
	filter := Filter{}
	var ok bool
	if filter.WorkFormatID, ok = queryID(c, "work_format_id"); !ok {
		return Filter{}, false
	}
	if filter.GroupID, ok = queryID(c, "group_id"); !ok {
		return Filter{}, false
	}
	if filter.LaneID, ok = queryID(c, "lane_id"); !ok {
		return Filter{}, false
	}
	if filter.WasHeld, ok = queryBool(c, "was_held"); !ok {
		return Filter{}, false
	}
	from, ok := parseDate(c, c.Query("from"), "from")
	if !ok {
		return Filter{}, false
	}
	to, ok := parseDate(c, c.Query("to"), "to")
	if !ok {
		return Filter{}, false
	}
	filter.From = from
	filter.To = to
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

func queryBool(c *gin.Context, name string) (*bool, bool) {
	raw := c.Query(name)
	if raw == "" {
		return nil, true
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return nil, false
	}
	return &value, true
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

func parseClock(c *gin.Context, raw, field string) (time.Time, bool) {
	if raw == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": field + " is required"})
		return time.Time{}, false
	}
	for _, layout := range []string{"15:04:05", "15:04"} {
		value, err := time.Parse(layout, raw)
		if err == nil {
			return value, true
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": field + " must be HH:MM"})
	return time.Time{}, false
}

func toSession(c *gin.Context, req SessionRequest, id int) (Session, bool) {
	sessionDate, ok := parseDate(c, req.SessionDate, "session_date")
	if !ok {
		return Session{}, false
	}
	if sessionDate.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_date is required"})
		return Session{}, false
	}
	startTime, ok := parseClock(c, req.StartTime, "start_time")
	if !ok {
		return Session{}, false
	}
	endTime, ok := parseClock(c, req.EndTime, "end_time")
	if !ok {
		return Session{}, false
	}
	return Session{
		ID:           id,
		WorkFormatID: req.WorkFormatID,
		GroupID:      req.GroupID,
		LaneID:       req.LaneID,
		SessionDate:  sessionDate,
		StartTime:    startTime,
		EndTime:      endTime,
		Topic:        req.Topic,
		WasHeld:      req.WasHeld,
		TrainerIDs:   req.TrainerIDs,
	}, true
}

func toSessionResponse(item Session) SessionResponse {
	trainers := item.TrainerIDs
	if trainers == nil {
		trainers = []int{}
	}
	return SessionResponse{
		ID:           item.ID,
		WorkFormatID: item.WorkFormatID,
		GroupID:      item.GroupID,
		LaneID:       item.LaneID,
		SessionDate:  item.SessionDate.Format("02.01.2006"),
		StartTime:    item.StartTime.Format("15:04"),
		EndTime:      item.EndTime.Format("15:04"),
		Topic:        item.Topic,
		WasHeld:      item.WasHeld,
		TrainerIDs:   trainers,
	}
}

func writeSessionError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
	case errors.Is(err, ErrWorkFormatNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "work format not found"})
	case errors.Is(err, ErrGroupNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
	case errors.Is(err, ErrLaneNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "lane not found"})
	case errors.Is(err, ErrTrainerNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "trainer not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
