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
	cleared := map[string]string{}
	for _, cookie := range rec.Result().Cookies() {
		if cookie.MaxAge < 0 && cookie.Value == "" && cookie.HttpOnly {
			cleared[cookie.Name] = cookie.Path
		}
	}
	if cleared[accessCookie] != accessPath || cleared[refreshCookie] != refreshPath {
		t.Fatalf("cleared=%v", cleared)
	}
}

func TestLoginRejectsBadCredentialsAndBlockedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	issuer := testIssuer(t, time.Minute, time.Hour)

	unauthorized := engineWithAuth(t, issuer, &stubUsers{authErr: user.ErrUnauthorized})
	rec := postJSON(unauthorized, "/api/v1/login", `{"login":"admin","password":"wrong"}`)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("bad password status=%d body=%s", rec.Code, rec.Body.String())
	}
	if len(rec.Result().Cookies()) != 0 {
		t.Fatalf("failed login set cookies: %v", rec.Result().Cookies())
	}

	blocked := engineWithAuth(t, issuer, &stubUsers{authErr: user.ErrBlocked})
	rec = postJSON(blocked, "/api/v1/login", `{"login":"admin","password":"password"}`)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("blocked status=%d body=%s", rec.Code, rec.Body.String())
	}

	malformed := engineWithAuth(t, issuer, &stubUsers{u: sampleUser()})
	rec = postJSON(malformed, "/api/v1/login", `{"login"`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("malformed status=%d", rec.Code)
	}
}

func TestRefreshRejectsMissingAndWrongToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	issuer := testIssuer(t, time.Minute, time.Hour)
	engine := engineWithAuth(t, issuer, &stubUsers{u: sampleUser()})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/refresh", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("missing cookie status=%d", rec.Code)
	}

	access, _, err := issuer.Issue(sampleUser().ID, sampleUser().Role)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	req = httptest.NewRequest(http.MethodPost, "/api/v1/refresh", nil)
	req.AddCookie(&http.Cookie{Name: refreshCookie, Value: access})
	rec = httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("access token as refresh status=%d", rec.Code)
	}
}

func TestRefreshRejectsBlockedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	issuer := testIssuer(t, time.Minute, time.Hour)
	blocked := sampleUser()
	blocked.Status = user.StatusBlocked
	engine := engineWithAuth(t, issuer, &stubUsers{u: blocked})

	_, refresh, err := issuer.Issue(blocked.ID, blocked.Role)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/refresh", nil)
	req.AddCookie(&http.Cookie{Name: refreshCookie, Value: refresh})
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestGuardLeavesLoginOpenAndProtectsMe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	issuer := testIssuer(t, time.Minute, time.Hour)
	store := &stubUsers{u: sampleUser()}
	engine := gin.New()
	api := engine.Group("/api/v1")
	api.Use(RequireUser(store, issuer))
	NewHandler(store, issuer).RegisterRoutes(api)

	rec := postJSON(engine, "/api/v1/login", `{"login":"admin","password":"password"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("login behind guard status=%d body=%s", rec.Code, rec.Body.String())
	}
	if store.gotLogin != "admin" || store.gotPassword != "password" {
		t.Fatalf("authenticate got login=%q password=%q", store.gotLogin, store.gotPassword)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	rec = httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("me without cookie status=%d", rec.Code)
	}

	_, refresh, err := issuer.Issue(sampleUser().ID, sampleUser().Role)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	req = httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.AddCookie(&http.Cookie{Name: accessCookie, Value: refresh})
	rec = httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("me with refresh token status=%d", rec.Code)
	}

	access, _, err := issuer.Issue(sampleUser().ID, sampleUser().Role)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	req = httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.AddCookie(&http.Cookie{Name: accessCookie, Value: access})
	rec = httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"login":"admin"`) {
		t.Fatalf("me status=%d body=%s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "password") {
		t.Fatalf("me leaked password fields: %s", rec.Body.String())
	}
}

func TestSecureCookieFlag(t *testing.T) {
	gin.SetMode(gin.TestMode)
	issuer, err := NewIssuer("test-secret-must-be-at-least-32-chars", time.Minute, time.Hour, true)
	if err != nil {
		t.Fatalf("NewIssuer: %v", err)
	}
	engine := engineWithAuth(t, issuer, &stubUsers{u: sampleUser()})
	rec := postJSON(engine, "/api/v1/login", `{"login":"admin","password":"password"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
	for _, cookie := range rec.Result().Cookies() {
		if !cookie.Secure {
			t.Fatalf("%s is not Secure", cookie.Name)
		}
	}
}

func TestRoleGuard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	issuer := testIssuer(t, time.Minute, time.Hour)
	store := &stubUsers{u: sampleUser()}
	store.u.Role = user.RoleClient
	engine := gin.New()
	api := engine.Group("/api/v1")
	api.Use(RequireUser(store, issuer))
	api.GET("/users/", func(c *gin.Context) { c.Status(http.StatusOK) })
	api.GET("/tariffs/", func(c *gin.Context) { c.Status(http.StatusOK) })
	api.GET("/clients/", func(c *gin.Context) { c.Status(http.StatusOK) })
	api.GET("/journal/", func(c *gin.Context) { c.Status(http.StatusOK) })
	api.GET("/journal/:id", func(c *gin.Context) { c.Status(http.StatusOK) })
	api.POST("/users/", func(c *gin.Context) { c.Status(http.StatusCreated) })
	api.POST("/enrollments/", func(c *gin.Context) { c.Status(http.StatusCreated) })
	api.POST("/standards/", func(c *gin.Context) { c.Status(http.StatusCreated) })
	api.PATCH("/enrollments/:service_id/:session_id", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	access, _, err := issuer.Issue(store.u.ID, store.u.Role)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	call := func(method, path string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, nil)
		req.AddCookie(&http.Cookie{Name: accessCookie, Value: access})
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)
		return rec
	}

	if rec := call(http.MethodGet, "/api/v1/users/"); rec.Code != http.StatusForbidden {
		t.Fatalf("client directory status=%d", rec.Code)
	}
	if rec := call(http.MethodGet, "/api/v1/tariffs/"); rec.Code != http.StatusForbidden {
		t.Fatalf("client tariffs status=%d", rec.Code)
	}
	if rec := call(http.MethodGet, "/api/v1/clients/"); rec.Code != http.StatusForbidden {
		t.Fatalf("client cards status=%d", rec.Code)
	}
	if rec := call(http.MethodPost, "/api/v1/enrollments/"); rec.Code != http.StatusCreated {
		t.Fatalf("client enroll status=%d body=%s", rec.Code, rec.Body.String())
	}
	if rec := call(http.MethodGet, "/api/v1/journal/"); rec.Code != http.StatusForbidden || !strings.Contains(rec.Body.String(), "forbidden") {
		t.Fatalf("client journal status=%d body=%s", rec.Code, rec.Body.String())
	}
	if rec := call(http.MethodGet, "/api/v1/journal/4"); rec.Code != http.StatusForbidden {
		t.Fatalf("client journal entry status=%d", rec.Code)
	}
	if rec := call(http.MethodPost, "/api/v1/users/"); rec.Code != http.StatusForbidden || !strings.Contains(rec.Body.String(), "forbidden") {
		t.Fatalf("client write status=%d body=%s", rec.Code, rec.Body.String())
	}

	store.u.Role = user.RoleTrainer
	if rec := call(http.MethodGet, "/api/v1/clients/"); rec.Code != http.StatusOK {
		t.Fatalf("trainer clients status=%d", rec.Code)
	}
	if rec := call(http.MethodGet, "/api/v1/users/"); rec.Code != http.StatusForbidden {
		t.Fatalf("trainer users status=%d", rec.Code)
	}
	if rec := call(http.MethodPost, "/api/v1/standards/"); rec.Code != http.StatusCreated {
		t.Fatalf("trainer standard status=%d body=%s", rec.Code, rec.Body.String())
	}
	if rec := call(http.MethodPatch, "/api/v1/enrollments/1/2"); rec.Code != http.StatusNoContent {
		t.Fatalf("trainer attendance status=%d body=%s", rec.Code, rec.Body.String())
	}
	if rec := call(http.MethodPost, "/api/v1/users/"); rec.Code != http.StatusForbidden {
		t.Fatalf("trainer user create status=%d", rec.Code)
	}
	if rec := call(http.MethodGet, "/api/v1/journal/"); rec.Code != http.StatusForbidden {
		t.Fatalf("trainer journal status=%d", rec.Code)
	}

	store.u.Role = user.RoleAdmin
	if rec := call(http.MethodPost, "/api/v1/users/"); rec.Code != http.StatusCreated {
		t.Fatalf("admin write status=%d body=%s", rec.Code, rec.Body.String())
	}
	if rec := call(http.MethodGet, "/api/v1/journal/"); rec.Code != http.StatusOK {
		t.Fatalf("admin journal status=%d", rec.Code)
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
	u           user.User
	authErr     error
	gotLogin    string
	gotPassword string
}

func (s *stubUsers) Authenticate(_ context.Context, login, password string) (user.User, error) {
	s.gotLogin = login
	s.gotPassword = password
	if s.authErr != nil {
		return user.User{}, s.authErr
	}
	return s.u, nil
}

func (s *stubUsers) GetByID(context.Context, int) (user.User, error) {
	return s.u, nil
}

func (s *stubUsers) Card(context.Context, int) (int, int, error) {
	return 0, 0, nil
}

func engineWithAuth(t *testing.T, issuer *Issuer, users *stubUsers) *gin.Engine {
	t.Helper()
	engine := gin.New()
	NewHandler(users, issuer).RegisterRoutes(engine.Group("/api/v1"))
	return engine
}

func postJSON(engine *gin.Engine, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	return rec
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
