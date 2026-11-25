package payment

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

type PaymentService interface {
	GetAll(ctx context.Context, filter Filter) ([]Payment, error)
	GetByID(ctx context.Context, ID int) (Payment, error)
	ClientIDByUser(ctx context.Context, userID int) (int, error)
	Create(ctx context.Context, payment Payment) (Payment, error)
	Update(ctx context.Context, payment Payment) error
	UpdateStatus(ctx context.Context, id int, status string) error
	Delete(ctx context.Context, id int) error
}

type Handler struct {
	service PaymentService
}

func NewHandler(service PaymentService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	payments := r.Group("/payments")
	{
		payments.GET("/", h.GetAll)
		payments.POST("/", h.Create)
		payments.GET("/:id", h.GetByID)
		payments.PUT("/:id", h.Update)
		payments.PATCH("/:id/status", h.UpdateStatus)
		payments.DELETE("/:id", h.Delete)
	}
}

type PaymentResponse struct {
	ID         int     `json:"id"`
	ServiceID  int     `json:"service_id"`
	Amount     float64 `json:"amount"`
	PaidOn     string  `json:"paid_on"`
	Method     string  `json:"method"`
	Purpose    *string `json:"purpose"`
	Status     string  `json:"status"`
	ClientName string  `json:"client_name"`
}

type PaymentRequest struct {
	ServiceID int     `json:"service_id"`
	Amount    float64 `json:"amount"`
	PaidOn    string  `json:"paid_on"`
	Method    string  `json:"method"`
	Purpose   string  `json:"purpose"`
	Status    string  `json:"status"`
}

type StatusRequest struct {
	Status string `json:"status"`
}

func (h *Handler) GetAll(c *gin.Context) {
	filter, ok := parseFilter(c)
	if !ok {
		return
	}
	filter, ok = h.ownPayments(c, filter)
	if !ok {
		return
	}

	payments, err := h.service.GetAll(c.Request.Context(), filter)
	if err != nil {
		writePaymentError(c, err)
		return
	}

	resp := make([]PaymentResponse, 0, len(payments))
	for _, payment := range payments {
		resp = append(resp, toPaymentResponse(payment))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	payment, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		writePaymentError(c, err)
		return
	}
	if !h.canSeePayment(c, payment.ClientID) {
		return
	}

	c.JSON(http.StatusOK, toPaymentResponse(payment))
}

func (h *Handler) ownPayments(c *gin.Context, filter Filter) (Filter, bool) {
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
		writePaymentError(c, err)
		return filter, false
	}
	if clientID == 0 {
		c.JSON(http.StatusOK, []PaymentResponse{})
		return filter, false
	}
	filter.ClientID = clientID
	return filter, true
}

func (h *Handler) canSeePayment(c *gin.Context, clientID int) bool {
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
		writePaymentError(c, err)
		return false
	}
	if own != clientID {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return false
	}
	return true
}

func (h *Handler) Create(c *gin.Context) {
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	payment, ok := toPayment(c, req, 0)
	if !ok {
		return
	}

	created, err := h.service.Create(c.Request.Context(), payment)
	if err != nil {
		writePaymentError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toPaymentResponse(created))
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	payment, ok := toPayment(c, req, id)
	if !ok {
		return
	}

	if err := h.service.Update(c.Request.Context(), payment); err != nil {
		writePaymentError(c, err)
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
		writePaymentError(c, err)
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
		writePaymentError(c, err)
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
	filter := Filter{
		Status: c.Query("status"),
		Method: c.Query("method"),
	}
	raw := c.Query("service_id")
	if raw == "" {
		return filter, true
	}
	id, err := strconv.Atoi(raw)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return Filter{}, false
	}
	filter.ServiceID = id
	return filter, true
}

func toPayment(c *gin.Context, req PaymentRequest, id int) (Payment, bool) {
	payment := Payment{
		ID:        id,
		ServiceID: req.ServiceID,
		Amount:    req.Amount,
		Method:    req.Method,
		Purpose:   optionalString(req.Purpose),
		Status:    req.Status,
	}
	if req.PaidOn == "" {
		return payment, true
	}
	paidOn, err := time.Parse("2006-01-02", req.PaidOn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "paid_on must be YYYY-MM-DD"})
		return Payment{}, false
	}
	payment.PaidOn = paidOn
	return payment, true
}

func toPaymentResponse(payment Payment) PaymentResponse {
	return PaymentResponse{
		ID:         payment.ID,
		ServiceID:  payment.ServiceID,
		Amount:     payment.Amount,
		PaidOn:     payment.PaidOn.Format("02.01.2006"),
		Method:     payment.Method,
		Purpose:    payment.Purpose,
		Status:     payment.Status,
		ClientName: payment.ClientName,
	}
}

func writePaymentError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
	case errors.Is(err, ErrServiceNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
