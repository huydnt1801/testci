package user

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/go-logr/logr"
	"github.com/huydnt1801/chuyende/internal/ent"
	entutil "github.com/huydnt1801/chuyende/internal/utils/ent"
	"github.com/huydnt1801/chuyende/pkg/log"
)

type Service interface {
	CreateUser(ctx context.Context, user *User, password string) (*User, time.Time, error)
	UpdateUser(ctx context.Context, params *UserUpdate) (*User, error)
	Authenticate(ctx context.Context, username, password string) (*User, error)
	Confirm(ctx context.Context, user *User, otp string) error
	FindUser(ctx context.Context, params *UserParams) (*User, error)
	SendConfirmationToken(ctx context.Context, user *User) (time.Time, error)
}

type ServiceImpl struct {
	logger         logr.Logger
	entClient      *ent.Client
	passwordHasher PasswordHasher
	otpSender      OTPService
}

func NewService(db *sql.DB) *ServiceImpl {
	client := entutil.NewClientFromDB(db)
	s := &ServiceImpl{
		logger:         log.ZapLogger(),
		entClient:      client,
		passwordHasher: &BcryptPasswordHasher{},
		otpSender:      NewOTPService(db),
	}
	return s
}

func (s *ServiceImpl) CreateUser(ctx context.Context, user *User, password string) (*User, time.Time, error) {
	repo := NewRepo(s.entClient)
	if !isValidPhoneNumber(user.PhoneNumber) {
		return nil, time.Time{}, InvalidPhoneError{}
	}
	if err := DefaultPasswordComplexity.ValidatePassword(password); err != nil {
		return nil, time.Time{}, err
	}
	hashedPw, err := s.passwordHasher.HashPassword(password)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed hashing password: %w", err)
	}
	user.Password = hashedPw
	u, err := repo.CreateUser(ctx, user)
	if err != nil {
		return nil, time.Time{}, err
	}
	var nextTime time.Time
	nextTime, err = s.otpSender.SendOTP(ctx, user.PhoneNumber)
	if err != nil {
		return nil, time.Time{}, err
	}
	return u, nextTime, nil
}

func (s *ServiceImpl) UpdateUser(ctx context.Context, params *UserUpdate) (*User, error) {
	repo := NewRepo(s.entClient)
	if password := params.Password; password != "" {
		if err := DefaultPasswordComplexity.ValidatePassword(password); err != nil {
			return nil, err
		}
		hashedPw, err := s.passwordHasher.HashPassword(password)
		if err != nil {
			return nil, fmt.Errorf("failed hashing password: %w", err)
		}
		params.Password = hashedPw
	}
	user, err := s.FindUser(ctx, &UserParams{
		ID: params.ID,
	})
	if err != nil {
		return nil, err
	}
	usr, err := repo.UpdateUser(ctx, user, params)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (s *ServiceImpl) Authenticate(ctx context.Context, phoneNumber, password string) (*User, error) {
	repo := NewRepo(s.entClient)

	if !isValidPhoneNumber(phoneNumber) {
		return nil, InvalidPhoneError{}
	}
	user, err := repo.FindUser(ctx, &UserParams{PhoneNumber: phoneNumber})
	if err != nil {
		if IsUserNotFound(err) {
			return nil, UserNotFoundError{}
		} else {
			return nil, fmt.Errorf("failed querying user from DB: %w", err)
		}

	}

	if !s.passwordHasher.CheckPasswordHash(password, user.Password) {
		return nil, InvalidPasswordError{}
	}

	if !user.Confirmed {
		return nil, ConfirmError{}
	}
	user.Password = ""
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *ServiceImpl) Confirm(ctx context.Context, usr *User, otp string) error {
	repo := NewRepo(s.entClient)
	ok, err := s.otpSender.VerifyOTP(ctx, usr, otp)
	if err != nil {
		return err
	}
	if !ok {
		return InvalidOTPError{}
	}

	confirmed := true
	_, err = repo.UpdateUser(ctx, usr, &UserUpdate{
		Confirmed: &confirmed,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceImpl) FindUser(ctx context.Context, params *UserParams) (*User, error) {
	repo := NewRepo(s.entClient)
	if params.PhoneNumber != "" {
		if !isValidPhoneNumber(params.PhoneNumber) {
			return nil, InvalidPhoneError{}
		}
	}
	user, err := repo.FindUser(ctx, params)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *ServiceImpl) SendConfirmationToken(ctx context.Context, user *User) (time.Time, error) {
	return s.otpSender.SendOTP(ctx, user.PhoneNumber)
}

func isValidPhoneNumber(phoneNumber string) bool {
	const PhoneNumberPattern = "^0[0-9]{9}$"
	if ok, _ := regexp.MatchString(PhoneNumberPattern, phoneNumber); !ok {
		return false
	}
	return true
}
