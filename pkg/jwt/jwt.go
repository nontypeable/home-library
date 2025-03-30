package jwt

import (
	"errors"
	"fmt"
	"home-library/pkg/config"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PayloadToken struct {
	UserID uuid.UUID
	jwt.StandardClaims
}

type JWT struct {
	cfg config.JWTConfig
}

func NewJWT(cfg config.JWTConfig) *JWT {
	return &JWT{cfg: cfg}
}

func (j *JWT) GenerateToken(payload PayloadToken) (token string, err error) {
	if j.cfg.Secret == "" {
		return "", errors.New("secret key is required")
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(j.cfg.Secret))
}

func (j *JWT) VerifyToken(c echo.Context, token string) error {
	if token == "" {
		return errors.New("token is required")
	}

	if len(token) < 7 || token[:7] != "Bearer " {
		return errors.New("invalid token format")
	}

	token = token[7:]

	newToken, err := jwt.ParseWithClaims(token, &PayloadToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.cfg.Secret), nil
	})
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	if !newToken.Valid {
		return errors.New("invalid token")
	}

	c.Set("jwt", token)

	return nil
}
