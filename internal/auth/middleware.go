package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"pool/internal/user"
)

const ctxUser = "auth_user"

type UserStore interface {
	Authenticate(ctx context.Context, login, password string) (user.User, error)
	GetByID(ctx context.Context, ID int) (user.User, error)
}

func RequireUser(users UserStore, issuer *Issuer) gin.HandlerFunc {
	public := map[string]struct{}{
		"/api/v1/login":   {},
		"/api/v1/refresh": {},
		"/api/v1/logout":  {},
	}
	return func(c *gin.Context) {
		if _, ok := public[c.Request.URL.Path]; ok {
			c.Next()
			return
		}

		cookie, err := c.Request.Cookie(accessCookie)
		if err != nil {
			abortUnauthorized(c)
			return
		}
		claims, err := issuer.ParseAccess(cookie.Value)
		if err != nil {
			abortUnauthorized(c)
			return
		}

		u, err := users.GetByID(c.Request.Context(), claims.UserID)
		if err != nil {
			if errors.Is(err, user.ErrNotFound) {
				abortUnauthorized(c)
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		if u.Status == user.StatusBlocked {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user is blocked"})
			return
		}
		c.Set(ctxUser, u)
		c.Next()
	}
}

func CurrentUser(c *gin.Context) (user.User, bool) {
	value, ok := c.Get(ctxUser)
	if !ok {
		return user.User{}, false
	}
	u, ok := value.(user.User)
	return u, ok
}

func abortUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}
