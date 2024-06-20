package jwt

import (
	"messenger.auth/internal/models/dto"
)

func (j *jwtImpl) refreshTokens(token *dto.RefreshToken) (*dto.Tokens, error) {
	userID, err := j.deleteToken(token.RefreshToken, "refresh")
	if err != nil {
		return nil, err
	}

	return j.generateTokens(userID)
}
