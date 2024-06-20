package jwt

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func (j *jwtImpl) createToken(claims jwt.MapClaims, timeNow time.Time, duration time.Duration, privateRSA *rsa.PrivateKey) (string, int, string, error) {
	exp := timeNow.Add(duration).Unix()
	uuid := uuid.New().String()

	claims["iat"] = timeNow.Unix()
	claims["exp"] = exp
	claims["uuid"] = uuid

	tknString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateRSA)

	return tknString, int(exp - timeNow.Unix()), uuid, err
}
