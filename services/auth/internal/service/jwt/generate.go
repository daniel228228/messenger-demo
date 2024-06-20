package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"

	"messenger.auth/internal/models/dto"
)

func (j *jwtImpl) generateTokens(userID string) (*dto.Tokens, error) {
	privateRSA, err := j.readRSAPrivateKeyFromString(j.config.JWT.RSAPrivate)
	if err != nil {
		panic("error read private key RSA from Config.JWT.RSAPrivate")
	}

	accessDuration, err := time.ParseDuration(j.config.JWT.Expired)
	if err != nil {
		panic("error parse Config.JWT.Expired")
	}

	refreshDuration, err := time.ParseDuration(j.config.JWT.RefreshExpired)
	if err != nil {
		panic("error parse Config.Auth.RefreshExpired")
	}

	timeNow := time.Now()
	tokens := &dto.Tokens{}

	var accessUUID, refreshUUID string

	accessClaims := make(jwt.MapClaims)
	accessClaims["user_id"] = userID
	accessClaims["token_type"] = "access"

	tokens.AccessToken, tokens.ExpiresIn, accessUUID, err = j.createToken(accessClaims, timeNow, accessDuration, privateRSA)
	if err != nil {
		return nil, err
	}

	refreshClaims := make(jwt.MapClaims)
	refreshClaims["user_id"] = userID
	refreshClaims["token_type"] = "refresh"

	tokens.RefreshToken, _, refreshUUID, err = j.createToken(refreshClaims, timeNow, refreshDuration, privateRSA)
	if err != nil {
		return nil, err
	}

	if err := j.cache.Set(accessUUID, userID, accessDuration); err != nil {
		j.log.Error().Str("from", "jwt.generateTokens (set accessToken to cache)").Err(errors.Wrapf(err, "cache: Set error")).Send()
		return nil, err
	}

	if err := j.cache.Set(refreshUUID, userID, refreshDuration); err != nil {
		j.log.Error().Str("from", "jwt.generateTokens (set refreshToken to cache)").Err(errors.Wrapf(err, "cache: Set error")).Send()
		return nil, err
	}

	return tokens, nil
}
