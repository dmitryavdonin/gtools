package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrSigingMethod    = "unexpected signing method: %v"
	ErrExpired         = "token is expired"
	ErrClaims          = "error parsing claims"
	ErrEmptySigningKey = "empty siging key"
)

type TokenManager interface {
	NewJWT(input JwtClaims, ttl time.Duration) (string, error)
	Parse(accessToken string) (result JwtClaims, err error)
	ParseFromGrpcContext(ctx context.Context) (result JwtClaims, err error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New(ErrEmptySigningKey)
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(input JwtClaims, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		input.UserID,
		input.UserName,
		input.UserRole,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) ParseFromGrpcContext(ctx context.Context) (result JwtClaims, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return result, errors.New("no metadata in context")
	}

	token := ""
	if len(md.Get("authorization")) > 0 {
		token = md.Get("authorization")[0]
	} else {
		return result, errors.New("no authorization header")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(token, prefix) {
		return result, status.Error(codes.Unauthenticated, `missing "Bearer " prefix in "Authorization" header`)
	}

	token = strings.Replace(token, "Bearer ", "", 1)

	return m.Parse(token)
}

func (m *Manager) Parse(accessToken string) (result JwtClaims, err error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(ErrSigingMethod, token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return result, errors.New(ErrClaims)
	}
	if err := claims.Valid(); err != nil {
		return result, err
	}

	result.UserID = claims.UserID
	result.UserName = claims.UserName
	result.UserRole = claims.UserRole

	return
}
