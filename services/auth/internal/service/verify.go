package service

import (
	"context"

	"github.com/pkg/errors"
	"messenger.api/go/api"

	"messenger.auth/internal/models/dto"
	"messenger.auth/internal/otp_service"
)

func (s *service) Verify(verify *dto.Verify) (*dto.Tokens, error) {
	if ok, err := s.otpService.CheckOtp(verify.Phone, verify.Password); err != nil {
		if errors.Is(err, otp_service.ErrInvalidPhoneNumber) {
			return nil, ErrInvalidPhoneNumber
		}

		return nil, ErrOtpServiceError
	} else if !ok {
		return nil, ErrIncorrectCode
	}

	user, err := s.usersService.CreateUser(context.Background(), &api.CreateUserRequest{
		Username: verify.Phone,
	})
	if err != nil {
		return nil, err
	}

	result, err := s.jwt.GenerateTokens(user.UserId)
	if err != nil {
		s.log.Error().Str("from", "VerifyUser (jwt.GenerateTokens)").Err(errors.Wrapf(err, "jwt.GenerateTokens error")).Send()
		return nil, ErrInternalError
	}

	return result, nil
}
