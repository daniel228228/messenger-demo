package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

func (j *jwtImpl) parseToken(token string) (jwt.MapClaims, error) {
	key, err := j.readRSAPublicKeyFromString(j.config.JWT.RSAPublic)
	if err != nil {
		return nil, err
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (any, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New(fmt.Sprintf("unexpected method: %s", jwtToken.Header["alg"]))
		}

		return key, nil
	})
	if err != nil {
		return nil, ErrExpiredToken
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
