package service

import (
	"github.com/pkg/errors"

	"messenger.auth/internal/models/dto"
	"messenger.auth/internal/otp_service"
)

func (s *service) InitVerify(initVerify *dto.InitVerify) error {
	if err := s.otpService.SendOtp(initVerify.Phone); err != nil {
		if errors.Is(err, otp_service.ErrInvalidPhoneNumber) {
			return ErrInvalidPhoneNumber
		}

		return ErrOtpServiceError
	}

	return nil
}
