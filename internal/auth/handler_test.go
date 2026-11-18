package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"pool/internal/user"
)

func TestLoginSetsHttpOnlyCookiesAndHidesTokens(t *testing.T) {
	gin.SetMode(gin.TestMode)
	issuer := testIssuer(t, time.Minute, time.Hour)
	handler := NewHandler(&stubUsers{u: sampleUser()}, issuer)
	engine := gin.New()
	handler.RegisterRoutes(engine.Group("/api/v1"))

	body := `{"login":"admin","password":"password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}

	cookies := map[string]*http.Cookie{}
	for _, cookie := range rec.Result().Cookies() {
		cookies[cookie.Name] = cookie
	}
	access := cookies[accessCookie]
	refresh := cookies[refreshCookie]
	if access == nil || refresh == nil {
		t.Fatalf("cookies=%v", cookies)
	}
	for _, cookie := range []*http.Cookie{access, refresh} {
		if !cookie.HttpOnly {
			t.Fatalf("%s is readable by script", cookie.Name)
		}
		if cookie.SameSite != http.SameSiteLaxMode {
			t.Fatalf("%s SameSite=%v", cookie.Name, cookie.SameSite)
		}
		if strings.Contains(rec.Body.String(), cookie.Value) {
			t.Fatalf("%s leaked into the response body", cookie.Name)
		}
	}
	if access.Path != accessPath || refresh.Path != refreshPath {
		t.Fatalf("paths access=%s refresh=%s", access.Path, refresh.Path)
	}

	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("json: %v", err)
	}
	if payload["login"] != "admin" || payload["password"] != nil || payload["password_hash"] != nil {
		t.Fatalf("payload=%v", payload)
	}
}

func TestRefreshRotatesCookies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	issuer := testIssuer(t, time.Minute, time.Hour)
	handler := NewHandler(&stubUsers{u: sampleUser()}, issuer)
	engine := gin.New()
	handler.RegisterRoutes(engine.Group("/api/v1"))

	_, refresh, err := issuer.Issue(sampleUser().ID, sampleUser().Role)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/refresh", nil)
	req.AddCookie(&http.Cookie{Name: refreshCookie, Value: refresh})
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	var next string
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == refreshCookie {
			next = cookie.Value
		}
	}
	if next == "" || next == refresh {
		t.Fatal("refresh cookie was not rotated")
	}
}

func TestLogoutClearsCookies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	issuer := testIssuer(t, time.Minute, time.Hour)
	handler := NewHandler(&stubUsers{u: sampleUser()}, issuer)
	engine := gin.New()
	handler.RegisterRoutes(engine.Group("/api/v1"))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status=%d", rec.Code)
	}
	cleared := map[string]bool{}
	for _, cookie := range rec.Result().Cookies() {
		if cookie.MaxAge < 0 && cookie.Value == "" && cookie.HttpOnly {
			cleared[cookie.Name] = true
		}
	}
	if !cleared[accessCookie] || !cleared[refreshCookie] {
		t.Fatalf("cleared=%v", cleared)
	}
}

func TestRequireUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	issuer := testIssuer(t, time.Minute, time.Hour)
	store := &stubUsers{u: sampleUser()}
	engine := gin.New()
	api := engine.Group("/api/v1")
	api.Use(RequireUser(store, issuer))
	api.GET("/me", func(c *gin.Context) {
		if _, ok := CurrentUser(c); !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("missing cookie status=%d", rec.Code)
	}

	access, _, err := issuer.Issue(sampleUser().ID, sampleUser().Role)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	req = httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.AddCookie(&http.Cookie{Name: accessCookie, Value: access})
	rec = httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		body, _ := io.ReadAll(rec.Body)
		t.Fatalf("status=%d body=%s", rec.Code, body)
	}

	store.u.Status = user.StatusBlocked
	req = httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.AddCookie(&http.Cookie{Name: accessCookie, Value: access})
	rec = httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("blocked status=%d", rec.Code)
	}
}

type stubUsers struct {
	u user.User
}

func (s *stubUsers) Authenticate(context.Context, string, string) (user.User, error) {
	return s.u, nil
}

func (s *stubUsers) GetByID(context.Context, int) (user.User, error) {
	return s.u, nil
}

func sampleUser() user.User {
	return user.User{
		ID:        7,
		Login:     "admin",
		LastName:  "Иванов",
		FirstName: "Иван",
		Role:      user.RoleAdmin,
		Status:    user.StatusActive,
	}
}
