package journal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"pool/internal/auth"
	"pool/internal/user"
)

type captureRepo struct {
	got Entry
}

func (c *captureRepo) GetAll(context.Context, Filter) ([]Entry, error) { return nil, nil }
func (c *captureRepo) GetByID(context.Context, int) (Entry, error)     { return Entry{}, ErrNotFound }
func (c *captureRepo) Create(_ context.Context, entry Entry) (Entry, error) {
	c.got = entry
	return entry, nil
}

func TestRecordWritesLoginAndSkipsReads(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &captureRepo{}
	svc := NewService(repo)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/v1/users/", nil)
	Record(svc)(ctx)
	if repo.got.Action != "" {
		t.Fatalf("logged a read: %+v", repo.got)
	}

	recorder = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/login", nil)
	auth.SetCurrentUser(ctx, user.User{ID: 7, Login: "admin"})
	Record(svc)(ctx)
	if repo.got.Action != ActionLogin || repo.got.Object != "users" || repo.got.Description != "вход в систему (admin)" {
		t.Fatalf("login log: %+v", repo.got)
	}

	repo.got = Entry{}
	recorder = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPut, "/api/v1/sessions/20412", nil)
	ctx.Set("request_body", []byte(`{"was_held":false,"work_format_id":2,"client_id":597}`))
	auth.SetCurrentUser(ctx, user.User{ID: 7})
	Record(svc)(ctx)
	want := "изменение занятия 20412: формат индивидуальное, статус занятия запланировано, клиент 597"
	if repo.got.Description != want {
		t.Fatalf("session log: got %q want %q", repo.got.Description, want)
	}

	repo.got = Entry{}
	recorder = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/services/", nil)
	ctx.Writer.WriteHeader(http.StatusBadRequest)
	auth.SetCurrentUser(ctx, user.User{ID: 7})
	Record(svc)(ctx)
	if repo.got.Action != ActionCreate || repo.got.Result != ResultError {
		t.Fatalf("failed create log: %+v", repo.got)
	}
	if !strings.Contains(repo.got.Description, "ошибка данных") {
		t.Fatalf("failed create description: %q", repo.got.Description)
	}

	repo.got = Entry{}
	recorder = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/login", nil)
	ctx.Writer.WriteHeader(http.StatusForbidden)
	auth.SetCurrentUser(ctx, user.User{ID: 7, Login: "blocked_user"})
	Record(svc)(ctx)
	if repo.got.Action != ActionLogin || repo.got.Result != ResultError {
		t.Fatalf("blocked login log: %+v", repo.got)
	}
	if repo.got.Description != "вход отклонён: учётка заблокирована (blocked_user)" {
		t.Fatalf("blocked login description: %q", repo.got.Description)
	}
}
