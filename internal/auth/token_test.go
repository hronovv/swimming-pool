package auth

import (
	"errors"
	"testing"
	"time"
)

func TestIssueAndParse(t *testing.T) {
	issuer := testIssuer(t, time.Minute, time.Hour)

	access, refresh, err := issuer.Issue(7, "admin")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if access == refresh {
		t.Fatal("access and refresh tokens must differ")
	}

	accessClaims, err := issuer.ParseAccess(access)
	if err != nil || accessClaims.UserID != 7 || accessClaims.Role != "admin" {
		t.Fatalf("access claims=%+v err=%v", accessClaims, err)
	}
	refreshClaims, err := issuer.ParseRefresh(refresh)
	if err != nil || refreshClaims.UserID != 7 {
		t.Fatalf("refresh claims=%+v err=%v", refreshClaims, err)
	}
	if _, err := issuer.ParseRefresh(access); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("access token used as refresh: %v", err)
	}
	if _, err := issuer.ParseAccess(refresh); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("refresh token used as access: %v", err)
	}
}

func TestParseRejectsExpiredToken(t *testing.T) {
	issuer := testIssuer(t, time.Minute, time.Hour)
	token, err := issuer.sign(3, "client", accessKind, time.Now().Add(-2*time.Hour), time.Now().Add(-time.Hour))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if _, err := issuer.ParseAccess(token); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("got %v, want ErrUnauthorized", err)
	}
}

func TestNewIssuerRejectsShortSecret(t *testing.T) {
	_, err := NewIssuer("short", time.Minute, time.Hour, false)
	if err == nil {
		t.Fatal("expected error")
	}
}

func testIssuer(t *testing.T, accessTTL, refreshTTL time.Duration) *Issuer {
	t.Helper()
	issuer, err := NewIssuer("test-secret-must-be-at-least-32-chars", accessTTL, refreshTTL, false)
	if err != nil {
		t.Fatalf("NewIssuer: %v", err)
	}
	return issuer
}
