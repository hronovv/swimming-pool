package journal

import (
	"context"
	"net/http"
	"net/http/httptest"
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
	auth.SetCurrentUser(ctx, user.User{ID: 7})
	Record(svc)(ctx)
	if repo.got.Action != ActionLogin || repo.got.Object != "users" || repo.got.Description != "login users, success" {
		t.Fatalf("login log: %+v", repo.got)
	}

	repo.got = Entry{}
	recorder = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/services/", nil)
	ctx.Writer.WriteHeader(http.StatusBadRequest)
	auth.SetCurrentUser(ctx, user.User{ID: 7})
	Record(svc)(ctx)
	if repo.got.Action != "" {
		t.Fatalf("logged a failed request: %+v", repo.got)
	}
}
