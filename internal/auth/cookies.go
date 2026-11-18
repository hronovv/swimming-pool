package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	accessCookie  = "access_token"
	refreshCookie = "refresh_token"

	accessPath  = "/api/v1"
	refreshPath = "/api/v1/refresh"
)

func (i *Issuer) SetCookies(c *gin.Context, access, refresh string) {
	http.SetCookie(c.Writer, sessionCookie(accessCookie, access, accessPath, int(i.accessTTL.Seconds()), i.secureCookie))
	http.SetCookie(c.Writer, sessionCookie(refreshCookie, refresh, refreshPath, int(i.refreshTTL.Seconds()), i.secureCookie))
}

func (i *Issuer) ClearCookies(c *gin.Context) {
	http.SetCookie(c.Writer, sessionCookie(accessCookie, "", accessPath, -1, i.secureCookie))
	http.SetCookie(c.Writer, sessionCookie(refreshCookie, "", refreshPath, -1, i.secureCookie))
}

func sessionCookie(name, value, path string, maxAge int, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}
