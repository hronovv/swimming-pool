package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"pool/internal/user"
)

type Handler struct {
	users  UserStore
	issuer *Issuer
}

func NewHandler(users UserStore, issuer *Issuer) *Handler {
	return &Handler{
		users:  users,
		issuer: issuer,
	}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.POST("/login", h.Login)
	r.POST("/refresh", h.Refresh)
	r.POST("/logout", h.Logout)
	r.GET("/me", h.Me)
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type userResponse struct {
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
	ClientID     *int    `json:"client_id,omitempty"`
	TrainerID    *int    `json:"trainer_id,omitempty"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	u, err := h.users.Authenticate(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if u.ID > 0 {
			SetCurrentUser(c, u)
		}
		writeAuthError(c, err)
		return
	}

	if err := h.issueCookies(c, u); err != nil {
		writeAuthError(c, err)
		return
	}
	SetCurrentUser(c, u)
	h.writeUser(c, http.StatusOK, u)
}

func (h *Handler) Refresh(c *gin.Context) {
	cookie, err := c.Request.Cookie(refreshCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	claims, err := h.issuer.ParseRefresh(cookie.Value)
	if err != nil {
		h.issuer.ClearCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	u, err := h.users.GetByID(c.Request.Context(), claims.UserID)
	if err != nil {
		h.issuer.ClearCookies(c)
		writeAuthError(c, err)
		return
	}
	if u.Status == user.StatusBlocked {
		h.issuer.ClearCookies(c)
		c.JSON(http.StatusForbidden, gin.H{"error": "user is blocked"})
		return
	}
	if err := h.issueCookies(c, u); err != nil {
		writeAuthError(c, err)
		return
	}
	h.writeUser(c, http.StatusOK, u)
}

func (h *Handler) Logout(c *gin.Context) {
	h.issuer.ClearCookies(c)
	c.Status(http.StatusNoContent)
}

func (h *Handler) Me(c *gin.Context) {
	u, ok := CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	h.writeUser(c, http.StatusOK, u)
}

func (h *Handler) writeUser(c *gin.Context, status int, u user.User) {
	clientID, trainerID, err := h.users.Card(c.Request.Context(), u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	resp := toUserResponse(u)
	if clientID > 0 {
		resp.ClientID = &clientID
	}
	if trainerID > 0 {
		resp.TrainerID = &trainerID
	}
	c.JSON(status, resp)
}

func (h *Handler) issueCookies(c *gin.Context, u user.User) error {
	access, refresh, err := h.issuer.Issue(u.ID, u.Role)
	if err != nil {
		return err
	}
	h.issuer.SetCookies(c, access, refresh)
	return nil
}

func toUserResponse(u user.User) userResponse {
	return userResponse{
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

func writeAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, user.ErrInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, user.ErrNotFound), errors.Is(err, user.ErrUnauthorized), errors.Is(err, ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	case errors.Is(err, user.ErrBlocked):
		c.JSON(http.StatusForbidden, gin.H{"error": "user is blocked"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
