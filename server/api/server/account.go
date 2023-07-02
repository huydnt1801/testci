package server

import (
	"crypto/hmac"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/huydnt1801/chuyende/api/server/middleware/auth"
	"github.com/huydnt1801/chuyende/internal/config"
	"github.com/huydnt1801/chuyende/internal/driver"
	"github.com/huydnt1801/chuyende/internal/user"
	"github.com/huydnt1801/chuyende/pkg/log"
	"github.com/labstack/echo/v4"
)

type AccountServer struct {
	logger    logr.Logger
	secretKey string

	userSvc   user.Service
	driverSvc driver.Service
}

func NewAccountServer(db *sql.DB) *AccountServer {
	cfg := config.MustParseConfig()
	svc := user.NewService(db)
	srv := &AccountServer{
		logger:    log.ZapLogger(),
		secretKey: cfg.SecretKey,
		userSvc:   svc,
		driverSvc: driver.NewService(db),
	}
	return srv
}

func (s *AccountServer) Register(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	data := &RegisterRequest{}
	if err := c.Bind(data); err != nil {
		return err
	}
	if err := c.Validate(data); err != nil {
		return err
	}
	if data.Password != data.Password2 {
		return echo.NewHTTPError(http.StatusBadRequest, "Mật khẩu không khớp")
	}
	_, nextOTPSend, err := s.userSvc.CreateUser(ctx, &user.User{
		PhoneNumber: data.PhoneNumber,
		FullName:    data.FullName,
	}, data.Password)
	if err != nil {
		return err
	}
	token := s.signConfirmInfo("reg", data.PhoneNumber, nextOTPSend)
	return c.JSON(http.StatusCreated, RegisterResponse{Code: http.StatusOK, Data: token})
}

type RegisterRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	FullName    string `json:"fullName" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Password2   string `json:"password2" validate:"required"`
}

type RegisterResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

func (s *AccountServer) ResendOTP(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	data := &ResendOTPRequest{}
	if err := c.Bind(data); err != nil {
		return err
	}
	if err := c.Validate(data); err != nil {
		return err
	}
	usr, err := s.userSvc.FindUser(ctx, &user.UserParams{PhoneNumber: data.PhoneNumber})
	if err != nil {
		return err
	}
	nextOTPSend, err := s.userSvc.SendConfirmationToken(ctx, usr)
	if err != nil {
		return err
	}
	token := s.signConfirmInfo("reg", data.PhoneNumber, nextOTPSend)
	return c.JSON(http.StatusCreated, RegisterResponse{Code: http.StatusOK, Data: token})
}

type ResendOTPRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required"`
}

type ResendOTPResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

func (s *AccountServer) RegisterConfirm(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	data := &RegisterConfirmRequest{}
	if err := c.Bind(data); err != nil {
		return err
	}
	if err := c.Validate(data); err != nil {
		return err
	}
	phoneNumber, nextOTPSend, err := s.parseConfirmInfo("reg", data.Token, 1*time.Hour)
	if err != nil {
		return err
	}
	if phoneNumber == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Sai thông tin")
	}
	usr, err := s.userSvc.FindUser(ctx, &user.UserParams{PhoneNumber: phoneNumber})
	if err != nil {
		return err
	}
	switch data.Type {
	case "submit-otp":
		err = s.userSvc.Confirm(ctx, usr, data.OTP)
		if err != nil {
			return err
		}
		auth.LoginUser(c, usr.ID, 0)
		return c.JSON(http.StatusOK, RegisterConfirmResponse{Code: http.StatusOK})
	case "resend-otp":
		if time.Now().Before(nextOTPSend) {
			return user.OTPIntervalError{}
		}
		nextOTPSend, err := s.userSvc.SendConfirmationToken(ctx, usr)
		if err != nil {
			return err
		}

		token := s.signConfirmInfo("reg", phoneNumber, nextOTPSend)
		return c.JSON(http.StatusCreated, RegisterConfirmResponse{Code: http.StatusOK, Data: token})
	default:
		return c.Redirect(http.StatusSeeOther, r.RequestURI)
	}
}

type RegisterConfirmRequest struct {
	Type  string `json:"type" validate:"required"`
	Token string `json:"token" validate:"required"`
	OTP   string `json:"otp" validate:"required"`
}

type RegisterConfirmResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

func (s *AccountServer) Logout(c echo.Context) error {
	auth.LogoutUser(c)
	return c.JSON(http.StatusOK, LogoutResponse{Code: http.StatusOK, Data: nil})
}

type LogoutResponse struct {
	Code int   `json:"code"`
	Data *bool `json:"data"`
}

func (s *AccountServer) Login(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	data := &LoginRequest{}
	if err := c.Bind(data); err != nil {
		return err
	}
	if err := c.Validate(data); err != nil {
		return err
	}
	usr, err := s.userSvc.Authenticate(ctx, data.PhoneNumber, data.Password)
	if err != nil {
		return err
	}
	if err := auth.LoginUser(c, usr.ID, 0); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, LoginResponse{Code: http.StatusOK, Data: usr})
}

type LoginRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Code int        `json:"code"`
	Data *user.User `json:"data"`
}

func (s *AccountServer) UpdateInfo(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	data := &UpdateInfoRequest{}
	if err := c.Bind(data); err != nil {
		return err
	}
	if err := c.Validate(data); err != nil {
		return err
	}
	authInfo, _ := auth.GetAuthInfo(c)
	if authInfo.UserID != data.UserID {
		return echo.NewHTTPError(http.StatusForbidden, "Bạn không có quyền chỉnh sửa")
	}
	usr, err := s.userSvc.UpdateUser(ctx, &user.UserUpdate{
		ID:       data.UserID,
		FullName: data.FullName,
		Password: data.Password,
		ImageURL: data.ImageURL,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, UpdateInfoResponse{Code: http.StatusOK, Data: usr})
}

type UpdateInfoRequest struct {
	UserID   int    `param:"userId" validate:"required,numeric"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
	ImageURL string `json:"imageUrl"`
}

type UpdateInfoResponse struct {
	Code int        `json:"code"`
	Data *user.User `json:"data"`
}

func (s *AccountServer) LoginDriver(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	data := &LoginDriverRequest{}
	if err := c.Bind(data); err != nil {
		return err
	}
	if err := c.Validate(data); err != nil {
		return err
	}
	driv, err := s.driverSvc.Authenticate(ctx, data.PhoneNumber, data.Password)
	if err != nil {
		return err
	}
	driv.Password = ""
	auth.LoginUser(c, 0, driv.ID)
	return c.JSON(http.StatusOK, LoginDriverResponse{Code: http.StatusOK, Data: driv})
}

type LoginDriverRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type LoginDriverResponse struct {
	Code int            `json:"code"`
	Data *driver.Driver `json:"data"`
}

func (s *AccountServer) CheckPhone(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	data := &CheckPhoneRequest{}
	if err := c.Bind(data); err != nil {
		return err
	}
	if err := c.Validate(data); err != nil {
		return err
	}
	usr, err := s.userSvc.FindUser(ctx, &user.UserParams{PhoneNumber: data.PhoneNumber})
	if user.IsUserNotFound(err) {
		return c.JSON(http.StatusOK, CheckPhoneResponse{Code: http.StatusOK, Data: nil})
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, CheckPhoneResponse{Code: http.StatusOK, Data: usr})
}

type CheckPhoneRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required"`
}

type CheckPhoneResponse struct {
	Code int        `json:"code"`
	Data *user.User `json:"data"`
}

func (s *AccountServer) signData(data []byte) []byte {
	mac := hmac.New(sha1.New, []byte(s.secretKey))
	mac.Write(data)
	return mac.Sum(nil)
}

func (s *AccountServer) signConfirmInfo(typ, phoneNumber string, nextOTPSend time.Time) string {
	regData := fmt.Sprintf("%s;%s;%d;%d", typ, phoneNumber, nextOTPSend.Unix(), time.Now().Unix())
	signature := s.signData([]byte(regData))
	encodeData := base64.StdEncoding.EncodeToString([]byte(regData))
	encodeSig := base64.StdEncoding.EncodeToString(signature)
	token := fmt.Sprintf("%s.%s", encodeData, encodeSig)
	return token
}

func (s *AccountServer) parseConfirmInfo(typ, hashed string, exp time.Duration) (string, time.Time, error) {
	ss := strings.Split(hashed, ".")
	if len(ss) != 2 {
		return "", time.Time{}, echo.NewHTTPError(http.StatusBadRequest, "Invalid token")
	}
	encodedData, sig := ss[0], ss[1]
	raw, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", time.Time{}, echo.NewHTTPError(http.StatusBadRequest, "Invalid token")
	}
	mac := hmac.New(sha1.New, []byte(s.secretKey))
	mac.Write(raw)
	expMac := mac.Sum(nil)
	tokenMac, _ := base64.StdEncoding.DecodeString(sig)
	if !hmac.Equal(tokenMac, expMac) {
		return "", time.Time{}, echo.NewHTTPError(http.StatusBadRequest, "Invalid token")
	}
	ss = strings.Split(string(raw), ";")
	if len(ss) != 4 || ss[0] != typ {
		return "", time.Time{}, echo.NewHTTPError(http.StatusBadRequest, "Invalid token")
	}
	nextOTPSendUnix, _ := strconv.ParseInt(ss[2], 10, 0)
	nextOTPSend := time.Unix(nextOTPSendUnix, 0)
	createdAtUnix, _ := strconv.ParseInt(ss[3], 10, 0)
	createdAt := time.Unix(createdAtUnix, 0)
	if createdAt.Add(exp).Before(time.Now()) {
		return "", time.Time{}, echo.NewHTTPError(http.StatusBadRequest, "Token is expired")
	}
	return ss[1], nextOTPSend, nil
}
