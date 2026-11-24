package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"pool/internal/user"
)

const ctxUser = "auth_user"

type UserStore interface {
	Authenticate(ctx context.Context, login, password string) (user.User, error)
	GetByID(ctx context.Context, ID int) (user.User, error)
	Card(ctx context.Context, userID int) (clientID int, trainerID int, err error)
}

func RequireUser(users UserStore, issuer *Issuer, after ...func(*gin.Context)) gin.HandlerFunc {
	var done func(*gin.Context)
	if len(after) > 0 {
		done = after[0]
	}
	public := map[string]struct{}{
		"/api/v1/login":   {},
		"/api/v1/refresh": {},
		"/api/v1/logout":  {},
	}
	return func(c *gin.Context) {
		if _, ok := public[c.Request.URL.Path]; ok {
			c.Next()
			finish(done, c)
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
		if !allowed(c.Request.Method, c.Request.URL.Path, u.Role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
		finish(done, c)
	}
}

func SetCurrentUser(c *gin.Context, u user.User) {
	c.Set(ctxUser, u)
}

func finish(done func(*gin.Context), c *gin.Context) {
	if done != nil {
		done(c)
	}
}

func allowed(method, path, role string) bool {
	if isJournal(path) || isDirectory(path) {
		return role == user.RoleAdmin
	}
	if isClients(path) {
		if role == user.RoleAdmin {
			return true
		}
		return role == user.RoleTrainer && readOnly(method)
	}
	if isPermits(path) {
		if role == user.RoleAdmin {
			return true
		}
		return role == user.RoleTrainer && readOnly(method)
	}
	if readOnly(method) {
		return true
	}
	if role == user.RoleAdmin {
		return true
	}
	if role == user.RoleTrainer {
		if method == http.MethodPost && path == "/api/v1/standards/" {
			return true
		}
		if method == http.MethodPatch && strings.HasPrefix(path, "/api/v1/enrollments/") {
			return true
		}
	}
	if role == user.RoleClient && method == http.MethodPost && path == "/api/v1/enrollments/" {
		return true
	}
	return false
}

func readOnly(method string) bool {
	return method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions
}

func isJournal(path string) bool {
	return path == "/api/v1/journal" || strings.HasPrefix(path, "/api/v1/journal/")
}

func isDirectory(path string) bool {
	for _, prefix := range []string{
		"/api/v1/users",
		"/api/v1/trainers",
		"/api/v1/tariffs",
		"/api/v1/groups",
		"/api/v1/lanes",
	} {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	return false
}

func isClients(path string) bool {
	return path == "/api/v1/clients" || strings.HasPrefix(path, "/api/v1/clients/")
}

func isPermits(path string) bool {
	return path == "/api/v1/permits" || strings.HasPrefix(path, "/api/v1/permits/")
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
