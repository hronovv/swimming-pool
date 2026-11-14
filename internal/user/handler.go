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

func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	c.JSON(http.StatusOK, user)
}
