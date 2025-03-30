package jwt

import (
	"errors"
	"home-library/pkg/config"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWT_GenerateToken(t *testing.T) {
	cfg := config.JWTConfig{
		Secret: "test-secret",
	}
	jwtService := NewJWT(cfg)
	userID := uuid.New()
	payload := PayloadToken{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	t.Run("successful token generation", func(t *testing.T) {
		token, err := jwtService.GenerateToken(payload)
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("empty secret", func(t *testing.T) {
		emptyCfg := config.JWTConfig{
			Secret: "",
		}
		emptyJWT := NewJWT(emptyCfg)
		token, err := emptyJWT.GenerateToken(payload)
		require.Error(t, err)
		require.Equal(t, "secret key is required", err.Error())
		require.Empty(t, token)
	})
}

func TestJWT_VerifyToken(t *testing.T) {
	cfg := config.JWTConfig{
		Secret: "test-secret",
	}
	jwtService := NewJWT(cfg)
	userID := uuid.New()
	payload := PayloadToken{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	validToken, err := jwtService.GenerateToken(payload)
	require.NoError(t, err)

	e := echo.New()
	c := e.NewContext(nil, nil)

	tests := []struct {
		name          string
		token         string
		expectedError error
	}{
		{
			name:          "valid token",
			token:         "Bearer " + validToken,
			expectedError: nil,
		},
		{
			name:          "invalid token format",
			token:         "invalid-token",
			expectedError: errors.New("invalid token format"),
		},
		{
			name:          "empty token",
			token:         "",
			expectedError: errors.New("token is required"),
		},
		{
			name:          "token without Bearer prefix",
			token:         validToken,
			expectedError: errors.New("invalid token format"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := jwtService.VerifyToken(c, tt.token)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			token, ok := c.Get("jwt").(string)
			assert.True(t, ok)
			assert.NotEmpty(t, token)
		})
	}
}
