package auth

import (
	"encoding/base64"
	"errors"
	"strings"
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

func TestNewIssuerRejectsBadTTL(t *testing.T) {
	secret := "test-secret-must-be-at-least-32-chars"
	cases := []struct {
		access, refresh time.Duration
	}{
		{0, time.Hour},
		{-time.Minute, time.Hour},
		{time.Hour, time.Hour},
		{time.Hour, time.Minute},
	}
	for _, tc := range cases {
		if _, err := NewIssuer(secret, tc.access, tc.refresh, false); err == nil {
			t.Fatalf("access=%s refresh=%s: expected error", tc.access, tc.refresh)
		}
	}
}

func TestParseRejectsForgedTokens(t *testing.T) {
	issuer := testIssuer(t, time.Minute, time.Hour)
	access, _, err := issuer.Issue(1, "client")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	parts := strings.Split(access, ".")
	sig := []byte(parts[2])
	if sig[0] == 'a' {
		sig[0] = 'b'
	} else {
		sig[0] = 'a'
	}
	parts[2] = string(sig)
	if _, err := issuer.ParseAccess(strings.Join(parts, ".")); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("tampered signature: %v", err)
	}

	other, err := NewIssuer("another-secret-must-be-32-characters", time.Minute, time.Hour, false)
	if err != nil {
		t.Fatalf("NewIssuer: %v", err)
	}
	if _, err := other.ParseAccess(access); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("foreign secret: %v", err)
	}

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"1","kind":"access","role":"admin","exp":4102444800}`))
	if _, err := issuer.ParseAccess(header + "." + payload + "."); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("alg none: %v", err)
	}
}

func TestIssueUsesDistinctTokenIDs(t *testing.T) {
	issuer := testIssuer(t, time.Minute, time.Hour)
	firstAccess, firstRefresh, err := issuer.Issue(1, "client")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	secondAccess, _, err := issuer.Issue(1, "client")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	first, err := issuer.ParseAccess(firstAccess)
	if err != nil {
		t.Fatalf("ParseAccess: %v", err)
	}
	second, err := issuer.ParseAccess(secondAccess)
	if err != nil {
		t.Fatalf("ParseAccess: %v", err)
	}
	refresh, err := issuer.ParseRefresh(firstRefresh)
	if err != nil {
		t.Fatalf("ParseRefresh: %v", err)
	}
	if first.ID == "" || first.ID == second.ID || first.ID == refresh.ID {
		t.Fatalf("jti first=%s second=%s refresh=%s", first.ID, second.ID, refresh.ID)
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
