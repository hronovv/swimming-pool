package user

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetUser(ctx context.Context, ID int) (User, error)
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
		users.GET("/:id", h.GetUser)
	}
}

type UserResponse struct {
	ID           int     `json:"id"`
	LastName     string  `json:"last_name"`
	FirstName    string  `json:"first_name"`
	Patronymic   *string `json:"patronymic"`
	Role         string  `json:"role"`
	Status       string  `json:"status"`
	RegisteredAt string  `json:"registered_at"`
	Phone        *string `json:"phone"`
	Email        *string `json:"email"`
}

func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	resp := UserResponse{
		ID:           user.ID,
		LastName:     user.LastName,
		FirstName:    user.FirstName,
		Patronymic:   user.Patronymic,
		Role:         user.Role,
		Status:       user.Status,
		RegisteredAt: user.RegisteredAt.Format("02.01.2006 15:04"),
		Phone:        user.Phone,
		Email:        user.Email,
	}

	c.JSON(http.StatusOK, resp)

}
