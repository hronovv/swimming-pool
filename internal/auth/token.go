package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessKind  = "access"
	refreshKind = "refresh"

	minSecretLen = 32
)

var ErrUnauthorized = errors.New("unauthorized")

type Claims struct {
	jwt.RegisteredClaims
	Role   string `json:"role"`
	Kind   string `json:"kind"`
	UserID int    `json:"-"`
}

type Issuer struct {
	secret       []byte
	accessTTL    time.Duration
	refreshTTL   time.Duration
	secureCookie bool
}

func NewIssuer(secret string, accessTTL, refreshTTL time.Duration, secureCookie bool) (*Issuer, error) {
	if len(secret) < minSecretLen {
		return nil, fmt.Errorf("JWT_SECRET must be at least %d characters", minSecretLen)
	}
	if accessTTL <= 0 || refreshTTL <= accessTTL {
		return nil, fmt.Errorf("refresh ttl must be longer than access ttl")
	}
	return &Issuer{
		secret:       []byte(secret),
		accessTTL:    accessTTL,
		refreshTTL:   refreshTTL,
		secureCookie: secureCookie,
	}, nil
}

func (i *Issuer) Issue(userID int, role string) (access, refresh string, err error) {
	if userID <= 0 {
		return "", "", fmt.Errorf("%w: user id is required", ErrUnauthorized)
	}
	now := time.Now()
	access, err = i.sign(userID, role, accessKind, now, now.Add(i.accessTTL))
	if err != nil {
		return "", "", err
	}
	refresh, err = i.sign(userID, role, refreshKind, now, now.Add(i.refreshTTL))
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (i *Issuer) ParseAccess(token string) (Claims, error) {
	return i.parse(token, accessKind)
}

func (i *Issuer) ParseRefresh(token string) (Claims, error) {
	return i.parse(token, refreshKind)
}

func (i *Issuer) sign(userID int, role, kind string, issued, expires time.Time) (string, error) {
	id, err := newTokenID()
	if err != nil {
		return "", err
	}
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id,
			Subject:   strconv.Itoa(userID),
			IssuedAt:  jwt.NewNumericDate(issued),
			ExpiresAt: jwt.NewNumericDate(expires),
		},
		Role: role,
		Kind: kind,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(i.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return token, nil
}

func (i *Issuer) parse(token, kind string) (Claims, error) {
	var claims Claims
	parsed, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnauthorized
		}
		return i.secret, nil
	})
	if err != nil || !parsed.Valid || claims.Kind != kind {
		return Claims{}, ErrUnauthorized
	}
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil || userID <= 0 {
		return Claims{}, ErrUnauthorized
	}
	claims.UserID = userID
	return claims, nil
}

func newTokenID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("token id: %w", err)
	}
	return hex.EncodeToString(buf), nil
}
