package otp_service

import (
	"context"

	"github.com/nyaruka/phonenumbers"
	"github.com/pkg/errors"

	"messenger.auth/pkg/log"

	"messenger.auth/pkg/config"
)

var Name = "OtpService"

type OtpService interface {
	SendOtp(to string) error
	CheckOtp(to, code string) (bool, error)
}

const (
	mockOtpCode = "1234"
)

var (
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
)

type mockOtpService struct {
	config *config.Config
	log    *log.Logger
}

func NewMockOtpService(config *config.Config, log *log.Logger) *mockOtpService {
	return &mockOtpService{
		config: config,
		log:    log,
	}
}

func (o *mockOtpService) Start(_ context.Context) error {
	o.log.Info().Bool("app", true).Str("component", Name).Str("state", "start").Send()

	return nil
}

func (o *mockOtpService) Stop(_ context.Context) error {
	o.log.Info().Bool("app", true).Str("component", Name).Str("state", "stop").Send()

	return nil
}

func (o *mockOtpService) SendOtp(to string) error {
	if !o.isValidPhoneNumber(to) {
		return ErrInvalidPhoneNumber
	}

	o.log.Debug().Str("from", "SendOtp").Str("phone", to).Msgf("Sent mock verification to %s", to)
	return nil
}

func (o *mockOtpService) CheckOtp(to, code string) (bool, error) {
	if !o.isValidPhoneNumber(to) {
		return false, ErrInvalidPhoneNumber
	} else {
		if code == mockOtpCode {
			o.log.Debug().Str("from", "CheckOtp").Msgf("Test user with phone %s and code %s successfully logged", to, code)
			return true, nil
		} else {
			o.log.Debug().Str("from", "CheckOtp").Msgf("Test user with phone %s and code %s unsuccessfully logged: incorrect code", to, code)
			return false, nil
		}
	}
}

func (o *mockOtpService) isValidPhoneNumber(phone string) bool {
	_, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return false
	}

	return true
}
