package server

import (
	"database/sql"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/huydnt1801/chuyende/api/server/middleware/auth"
	entTrip "github.com/huydnt1801/chuyende/internal/ent/trip"
	"github.com/huydnt1801/chuyende/internal/trip"
	"github.com/huydnt1801/chuyende/pkg/log"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
)

type TripServer struct {
	logger logr.Logger
	svc    trip.Service
}

func NewTripServer(db *sql.DB) *TripServer {
	svc := trip.NewService(db)
	srv := &TripServer{
		logger: log.ZapLogger(),
		svc:    svc,
	}
	return srv
}

type GetPriceTripRequest struct {
	Distance float64 `query:"distance" validate:"required,numeric"`
}
type TypeResponse struct {
	Motor float64 `json:"motor"`
	Car   float64 `json:"car"`
}

type GetPriceTripResponse struct {
	Code int           `json:"code"`
	Data *TypeResponse `json:"data"`
}

func (s *TripServer) GetPriceTrip(c echo.Context) error {
	data := &GetPriceTripRequest{}
	if err := c.Bind(data); err != nil {
		return err
	}
	if err := c.Validate(data); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GetPriceTripResponse{Code: http.StatusOK, Data: &TypeResponse{Motor: data.Distance * 20 * 1000, Car: data.Distance * 30 * 1000}})
}

type ListTripRequest struct {
	TripID *int    `query:"tripId"`
	Status *string `query:"status"`
	Rate   *int    `query:"rate"`
}

type ListTripResponse struct {
	Code int          `json:"code"`
	Data []*trip.Trip `json:"data"`
}

func (s *TripServer) ListTrip(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	authInfo, _ := auth.GetAuthInfo(c)
	data := &ListTripRequest{}
	if err := c.Bind(data); err != nil {
		return err
	}
	if err := c.Validate(data); err != nil {
		return err
	}

	params := &trip.TripParams{}
	err := mapstructure.Decode(data, params)
	if err != nil {
		return err
	}
	if authInfo.UserID != 0 {
		params.UserID = &authInfo.UserID
	}
	if authInfo.DriverID != 0 {
		waiting := entTrip.StatusWaiting
		params = &trip.TripParams{
			Status: &waiting,
		}
	}
	trip, err := s.svc.FindTrip(ctx, params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ListTripResponse{Code: http.StatusOK, Data: trip})
}

type CreateTripRequest struct {
	StartX        float64 `json:"startX" validate:"required,numeric"`
	StartY        float64 `json:"startY" validate:"required,numeric"`
	StartLocation string  `json:"startLocation" validate:"required"`
	EndX          float64 `json:"endX" validate:"required,numeric"`
	EndY          float64 `json:"endY" validate:"required,numeric"`
	EndLocation   string  `json:"endLocation" validate:"required"`
	Type          string  `json:"type" validate:"required"`
	Distance      float64 `json:"distance" validate:"required,numeric"`
}

type CreateTripResponse struct {
	Code int        `json:"code"`
	Data *trip.Trip `json:"data"`
}

func (s *TripServer) CreateTrip(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	authInfo, _ := auth.GetAuthInfo(c)
	if authInfo.UserID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Yêu cầu đăng nhập")
	}
	data := &CreateTripRequest{}
	if err := c.Bind(data); err != nil {
		return err
	}
	if err := c.Validate(data); err != nil {
		return err
	}

	createParams := &trip.Trip{}
	err := mapstructure.Decode(data, createParams)
	if err != nil {
		return err
	}
	createParams.UserID = authInfo.UserID
	trip, err := s.svc.CreateTrip(ctx, createParams)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, CreateTripResponse{Code: http.StatusOK, Data: trip})
}

type UpdateStatusTripRequest struct {
	TripID int    `param:"tripId" validate:"required,numeric"`
	Status string `json:"status" validate:"required"`
}

type UpdateStatusTripResponse struct {
	Code int        `json:"code"`
	Data *trip.Trip `json:"data"`
}

func (s *TripServer) UpdateStatusTrip(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	authInfo, _ := auth.GetAuthInfo(c)
	data := &UpdateStatusTripRequest{}
	if err := c.Bind(data); err != nil {
		return err
	}
	if err := c.Validate(data); err != nil {
		return err
	}

	updateParams := &trip.TripUpdate{}
	err := mapstructure.Decode(data, updateParams)
	if err != nil {
		return err
	}
	if updateParams.Status == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Trạng thái không hợp lệ")
	}
	tripParams := &trip.TripParams{
		TripID: &data.TripID,
	}
	if authInfo.UserID != 0 {
		if *updateParams.Status != entTrip.StatusCancel {
			return echo.NewHTTPError(http.StatusBadRequest, "Trạng thái không hợp lệ")
		}
		tripParams.UserID = &authInfo.UserID
	} else {
		if authInfo.DriverID == 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, "Yêu cầu đăng nhập")
		}
	}
	tripFounds, err := s.svc.FindTrip(ctx, tripParams)
	if err != nil {
		return err
	}
	if len(tripFounds) != 1 {
		return echo.NewHTTPError(http.StatusForbidden, "Bạn không có quyền chỉnh sửa")
	}
	if tripFounds[0].Status == entTrip.StatusCancel {
		return echo.NewHTTPError(http.StatusBadRequest, "Chuyến xe đã bị hủy")
	}
	if *updateParams.Status == entTrip.StatusCancel && tripFounds[0].Status == entTrip.StatusDone {
		return echo.NewHTTPError(http.StatusBadRequest, "Chuyến xe đã hoàn thành")
	}
	if authInfo.DriverID != 0 {
		if tripFounds[0].DriveID != 0 && authInfo.DriverID != tripFounds[0].DriveID {
			return echo.NewHTTPError(http.StatusForbidden, "Bạn không có quyền chỉnh sửa")
		}
		if *updateParams.Status != entTrip.StatusDone && *updateParams.Status != entTrip.StatusAccept {
			return echo.NewHTTPError(http.StatusBadRequest, "Trạng thái không hợp lệ")
		}
		updateParams.DriveID = &authInfo.DriverID
	}
	trip, err := s.svc.UpdateTrip(ctx, data.TripID, updateParams)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, UpdateStatusTripResponse{Code: http.StatusOK, Data: trip})
}

type RateTripRequest struct {
	TripID int `param:"tripId" validate:"required,numeric"`
	Rate   int `json:"rate" validate:"required,numeric"`
}

type RateTripResponse struct {
	Code int        `json:"code"`
	Data *trip.Trip `json:"data"`
}

func (s *TripServer) RateTrip(c echo.Context) error {
	r := c.Request()
	ctx := r.Context()
	authInfo, _ := auth.GetAuthInfo(c)
	if authInfo.UserID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Yêu cầu đăng nhập")
	}
	data := &RateTripRequest{}
	if err := c.Bind(data); err != nil {
		return err
	}
	if err := c.Validate(data); err != nil {
		return err
	}

	updateParams := &trip.TripUpdate{}
	err := mapstructure.Decode(data, updateParams)
	if err != nil {
		return err
	}
	tripFounds, err := s.svc.FindTrip(ctx, &trip.TripParams{
		TripID: &data.TripID,
		UserID: &authInfo.UserID,
	})
	if err != nil {
		return err
	}
	if len(tripFounds) != 1 {
		return echo.NewHTTPError(http.StatusForbidden, "Bạn không có quyền đánh giá")
	}
	trip, err := s.svc.UpdateTrip(ctx, data.TripID, updateParams)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, RateTripResponse{Code: http.StatusOK, Data: trip})
}
