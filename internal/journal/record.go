package journal

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"pool/internal/auth"
)

func Record(svc *Service) func(*gin.Context) {
	return func(c *gin.Context) {
		record(svc, c)
	}
}

func record(svc *Service, c *gin.Context) {
	if c.Writer.Status() >= http.StatusBadRequest {
		return
	}
	method := c.Request.Method
	if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
		return
	}
	path := c.Request.URL.Path
	if path == "/api/v1/refresh" || path == "/api/v1/logout" {
		return
	}
	u, ok := auth.CurrentUser(c)
	if !ok {
		return
	}
	action := ActionUpdate
	switch method {
	case http.MethodPost:
		action = ActionCreate
	case http.MethodDelete:
		action = ActionDelete
	}
	object := resource(path)
	if path == "/api/v1/login" {
		action = ActionLogin
		object = "users"
	}
	_, _ = svc.Create(c.Request.Context(), Entry{
		UserID:      u.ID,
		Action:      action,
		Object:      object,
		Description: action + " " + object + ", success",
		Result:      ResultSuccess,
		IPAddress:   c.ClientIP(),
	})
}

func resource(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 && parts[2] != "" {
		return parts[2]
	}
	return "api"
}
