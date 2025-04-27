package jwt

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type TokenData struct {
	ID int64
}

var (
	ErrTokenExpired = errors.New("token is expired")
	ErrTokenInvalid = errors.New("token is invalid")
)

// ParseToken throws ErrTokenExpired and ErrTokenInvalid
func ParseToken(token string) (*TokenData, error) {
	cl, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}

		return nil, ErrTokenInvalid
	}

	mp := cl.Claims.(jwt.MapClaims)

	ids, ok := mp["uid"]
	if !ok {
		return nil, ErrTokenInvalid
	}

	id, ok := ids.(int64)
	if !ok {
		return nil, ErrTokenInvalid
	}

	return &TokenData{ID: id}, nil
}
