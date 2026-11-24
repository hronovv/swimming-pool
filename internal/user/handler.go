package user

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetAll(ctx context.Context, filter Filter) ([]User, error)
	GetByID(ctx context.Context, ID int) (User, error)
	Create(ctx context.Context, in NewUser) (User, error)
	Update(ctx context.Context, id int, in UserUpdate) error
	UpdateStatus(ctx context.Context, id int, status string) error
	UpdatePassword(ctx context.Context, id int, password string) error
}

type Handler struct {
	service UserService
}

func NewHandler(service UserService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	users := r.Group("/users")
	{
		users.GET("/", h.GetAll)
		users.POST("/", h.Create)
		users.GET("/:id", h.GetByID)
		users.PUT("/:id", h.Update)
		users.PATCH("/:id/status", h.UpdateStatus)
		users.PATCH("/:id/password", h.UpdatePassword)
	}
}

type UserResponse struct {
	ID           int     `json:"id"`
	Login        string  `json:"login"`
	LastName     string  `json:"last_name"`
	FirstName    string  `json:"first_name"`
	Patronymic   *string `json:"patronymic"`
	Role         string  `json:"role"`
	Status       string  `json:"status"`
	RegisteredAt string  `json:"registered_at"`
	Phone        *string `json:"phone"`
	Email        *string `json:"email"`
}

type UserRequest struct {
	Login              string `json:"login"`
	Password           string `json:"password"`
	LastName           string `json:"last_name"`
	FirstName          string `json:"first_name"`
	Patronymic         string `json:"patronymic"`
	Role               string `json:"role"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
	Passport           string `json:"passport"`
	BirthDate          string `json:"birth_date"`
	Address            string `json:"address"`
	MedicalCertificate bool   `json:"medical_certificate"`
	Qualification      string `json:"qualification"`
	Rank               string `json:"rank"`
	ExperienceYears    int    `json:"experience_years"`
	Specialization     string `json:"specialization"`
}

type StatusRequest struct {
	Status string `json:"status"`
}

type PasswordRequest struct {
	Password string `json:"password"`
}

func (h *Handler) GetAll(c *gin.Context) {
	users, err := h.service.GetAll(c.Request.Context(), Filter{
		Role:   c.Query("role"),
		Status: c.Query("status"),
	})
	if err != nil {
		writeUserError(c, err)
		return
	}

	resp := make([]UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, toUserResponse(u))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	u, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		writeUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, toUserResponse(u))
}

func (h *Handler) Create(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var birthDate time.Time
	if strings.TrimSpace(req.BirthDate) != "" {
		parsed, err := time.Parse("2006-01-02", strings.TrimSpace(req.BirthDate))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user: birth_date must be YYYY-MM-DD"})
			return
		}
		birthDate = parsed
	}
	created, err := h.service.Create(c.Request.Context(), NewUser{
		Login:      req.Login,
		Password:   req.Password,
		LastName:   req.LastName,
		FirstName:  req.FirstName,
		Patronymic: req.Patronymic,
		Role:       req.Role,
		Phone:      req.Phone,
		Email:      req.Email,
		Profile: Profile{
			Passport:           req.Passport,
			BirthDate:          birthDate,
			Address:            req.Address,
			MedicalCertificate: req.MedicalCertificate,
			Qualification:      req.Qualification,
			Rank:               req.Rank,
			ExperienceYears:    req.ExperienceYears,
			Specialization:     req.Specialization,
		},
	})
	if err != nil {
		writeUserError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toUserResponse(created))
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var birthDate time.Time
	if strings.TrimSpace(req.BirthDate) != "" {
		parsed, err := time.Parse("2006-01-02", strings.TrimSpace(req.BirthDate))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user: birth_date must be YYYY-MM-DD"})
			return
		}
		birthDate = parsed
	}
	err := h.service.Update(c.Request.Context(), id, UserUpdate{
		Login:      req.Login,
		LastName:   req.LastName,
		FirstName:  req.FirstName,
		Patronymic: req.Patronymic,
		Role:       req.Role,
		Phone:      req.Phone,
		Email:      req.Email,
		Profile: Profile{
			Passport:           req.Passport,
			BirthDate:          birthDate,
			Address:            req.Address,
			MedicalCertificate: req.MedicalCertificate,
			Qualification:      req.Qualification,
			Rank:               req.Rank,
			ExperienceYears:    req.ExperienceYears,
			Specialization:     req.Specialization,
		},
	})
	if err != nil {
		writeUserError(c, err)
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
		writeUserError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdatePassword(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req PasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.service.UpdatePassword(c.Request.Context(), id, req.Password); err != nil {
		writeUserError(c, err)
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

func toUserResponse(u User) UserResponse {
	return UserResponse{
		ID:           u.ID,
		Login:        u.Login,
		LastName:     u.LastName,
		FirstName:    u.FirstName,
		Patronymic:   u.Patronymic,
		Role:         u.Role,
		Status:       u.Status,
		RegisteredAt: u.RegisteredAt.Format("02.01.2006"),
		Phone:        u.Phone,
		Email:        u.Email,
	}
}

func writeUserError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	case errors.Is(err, ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
	case errors.Is(err, ErrPassportTaken):
		c.JSON(http.StatusConflict, gin.H{"error": "passport already exists"})
	case errors.Is(err, ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	case errors.Is(err, ErrBlocked):
		c.JSON(http.StatusForbidden, gin.H{"error": "user is blocked"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
