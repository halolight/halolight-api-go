package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Claims struct {
	UserID    string    `json:"userId"`
	TokenType TokenType `json:"type"`
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"` // seconds
}

// GenerateAccessToken generates a new JWT access token (7 days default)
func GenerateAccessToken(userID, secret string) (string, error) {
	now := time.Now()
	expirationTime := now.Add(7 * 24 * time.Hour) // 7 days

	claims := Claims{
		UserID:    userID,
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken generates a new JWT refresh token (30 days default)
func GenerateRefreshToken(userID, secret string) (string, error) {
	now := time.Now()
	expirationTime := now.Add(30 * 24 * time.Hour) // 30 days

	claims := Claims{
		UserID:    userID,
		TokenType: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateTokenPair generates both access and refresh tokens
func GenerateTokenPair(userID, secret string) (*TokenPair, error) {
	accessToken, err := GenerateAccessToken(userID, secret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken(userID, secret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7 * 24 * 60 * 60, // 7 days in seconds
	}, nil
}

// ParseToken parses and validates a JWT token
func ParseToken(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// ValidateAccessToken validates an access token specifically
func ValidateAccessToken(tokenStr, secret string) (*Claims, error) {
	claims, err := ParseToken(tokenStr, secret)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != TokenTypeAccess {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token specifically
func ValidateRefreshToken(tokenStr, secret string) (*Claims, error) {
	claims, err := ParseToken(tokenStr, secret)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != TokenTypeRefresh {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

// GetTokenExpiration calculates token expiration time
func GetTokenExpiration(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}

// GetRefreshTokenExpiration returns 30 days from now
func GetRefreshTokenExpiration() time.Time {
	return GetTokenExpiration(30 * 24 * time.Hour)
}
