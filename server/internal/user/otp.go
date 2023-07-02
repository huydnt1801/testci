package user

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/huydnt1801/chuyende/internal/ent"
	"github.com/huydnt1801/chuyende/internal/ent/otp"
)

var TimeNow = time.Now

type OTPService interface {
	SendOTP(ctx context.Context, phoneNumber string) (time.Time, error)
	VerifyOTP(ctx context.Context, user *User, otp string) (bool, error)
}

func NewOTPService(db *sql.DB) *OTPServiceImpl {
	drv := entsql.OpenDB("mysql", db)
	client := ent.NewClient(ent.Driver(drv))
	return &OTPServiceImpl{
		client:              client,
		otpNextSendInterval: 1 * time.Minute,
		optExpireIn:         3 * time.Minute,
	}
}

type OTPServiceImpl struct {
	client              *ent.Client
	otpNextSendInterval time.Duration
	optExpireIn         time.Duration
}

func (s *OTPServiceImpl) SendOTP(ctx context.Context, phoneNumber string) (time.Time, error) {
	err := s.checkTimeInterval(ctx, phoneNumber)
	if err != nil {
		return time.Time{}, err
	}

	nextOTPSend := TimeNow().Add(s.otpNextSendInterval)

	otp, err := generateOTP(6)
	if err != nil {
		return time.Time{}, err
	}

	tx, err := s.client.Tx(ctx)
	if err != nil {
		return time.Time{}, fmt.Errorf("starting a transaction: %w", err)
	}
	if err := tx.Otp.Create().
		SetPhoneNumber(phoneNumber).
		SetOtp(otp).
		SetCreatedAt(TimeNow()).
		OnConflict().
		UpdateNewValues().
		Exec(ctx); err != nil {
		return time.Time{}, rollback(tx, err)
	}
	// send otp
	tx.Commit()
	return nextOTPSend, nil
}

func (s *OTPServiceImpl) checkTimeInterval(ctx context.Context, phoneNumber string) error {
	q := s.client.Otp.Query()
	q.Where(otp.PhoneNumber(phoneNumber))
	otp, _ := q.Only(ctx)
	if otp == nil {
		return nil
	}
	now := time.Now()
	nextOTPSend := otp.CreatedAt.Add(s.otpNextSendInterval)
	if now.Before(nextOTPSend) {
		return OTPIntervalError{}
	}
	return nil
}

func (s *OTPServiceImpl) VerifyOTP(ctx context.Context, user *User, code string) (bool, error) {
	q := s.client.Otp.Query().Where(otp.Otp(code))
	q = q.Where(otp.PhoneNumberEQ(user.PhoneNumber))
	rec, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	expired := rec.CreatedAt.Add(s.optExpireIn).Before(TimeNow())

	if err := s.client.Otp.DeleteOne(rec).Exec(ctx); err != nil {
		return false, err
	}
	return !expired, nil
}

func generateOTP(n int) (string, error) {
	const letters = "0123456789"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

// rollback calls to tx.Rollback and wraps the given error
// with the rollback error if occurred.
func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}
	return err
}
