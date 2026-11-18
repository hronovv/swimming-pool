package journal

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type JournalService interface {
	GetAll(ctx context.Context, filter Filter) ([]Entry, error)
	GetByID(ctx context.Context, ID int) (Entry, error)
	Create(ctx context.Context, entry Entry) (Entry, error)
}

type Handler struct {
	service JournalService
}

func NewHandler(service JournalService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	journal := r.Group("/journal")
	{
		journal.GET("/", h.GetAll)
		journal.POST("/", h.Create)
		journal.GET("/:id", h.GetByID)
	}
}

type EntryResponse struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Login       string `json:"login"`
	OccurredAt  string `json:"occurred_at"`
	Action      string `json:"action"`
	Object      string `json:"object"`
	Description string `json:"description"`
	Result      string `json:"result"`
	IPAddress   string `json:"ip_address"`
}

type EntryRequest struct {
	UserID      int    `json:"user_id"`
	Action      string `json:"action"`
	Object      string `json:"object"`
	Description string `json:"description"`
	Result      string `json:"result"`
	IPAddress   string `json:"ip_address"`
}

func (h *Handler) GetAll(c *gin.Context) {
	filter, ok := parseFilter(c)
	if !ok {
		return
	}

	entries, err := h.service.GetAll(c.Request.Context(), filter)
	if err != nil {
		writeJournalError(c, err)
		return
	}

	resp := make([]EntryResponse, 0, len(entries))
	for _, entry := range entries {
		resp = append(resp, toEntryResponse(entry))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	entry, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		writeJournalError(c, err)
		return
	}

	c.JSON(http.StatusOK, toEntryResponse(entry))
}

func (h *Handler) Create(c *gin.Context) {
	var req EntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	created, err := h.service.Create(c.Request.Context(), Entry{
		UserID:      req.UserID,
		Action:      req.Action,
		Object:      req.Object,
		Description: req.Description,
		Result:      req.Result,
		IPAddress:   req.IPAddress,
	})
	if err != nil {
		writeJournalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toEntryResponse(created))
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
	filter := Filter{
		Action: c.Query("action"),
		Result: c.Query("result"),
		Query:  strings.TrimSpace(c.Query("q")),
	}

	rawUserID := c.Query("user_id")
	if rawUserID != "" {
		id, err := strconv.Atoi(rawUserID)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return Filter{}, false
		}
		filter.UserID = id
	}

	from, ok := parseDate(c, c.Query("from"), "from")
	if !ok {
		return Filter{}, false
	}
	filter.From = from

	to, ok := parseDate(c, c.Query("to"), "to")
	if !ok {
		return Filter{}, false
	}
	if !to.IsZero() {
		filter.To = to.Add(24 * time.Hour)
	}

	return filter, true
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

func toEntryResponse(entry Entry) EntryResponse {
	return EntryResponse{
		ID:          entry.ID,
		UserID:      entry.UserID,
		Login:       entry.Login,
		OccurredAt:  entry.OccurredAt.Format("02.01.2006 15:04"),
		Action:      entry.Action,
		Object:      entry.Object,
		Description: entry.Description,
		Result:      entry.Result,
		IPAddress:   entry.IPAddress,
	}
}

func writeJournalError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "journal entry not found"})
	case errors.Is(err, ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
