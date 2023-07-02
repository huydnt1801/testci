//go:build integration
// +build integration

package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/huydnt1801/chuyende/api/server/middleware/auth"
	"github.com/huydnt1801/chuyende/internal/ent"
	entTrip "github.com/huydnt1801/chuyende/internal/ent/trip"
	"github.com/huydnt1801/chuyende/internal/trip"
	"github.com/huydnt1801/chuyende/test"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
)

const (
	TripEndpoint string = "/api/v1/trips"
)

var container testcontainers.Container

func TestMain(m *testing.M) {
	// Setup
	var cleanUp func()
	container, cleanUp = test.MysqlContainer(context.Background())

	// Run test case
	exitCode := m.Run()

	// After test
	cleanUp()

	os.Exit(exitCode)
}

type TestGetPriceTripInfo struct {
	authInfo    *auth.AuthInfo
	queryParams map[string]string
	output      *GetPriceTripResponse
}

func TestGetPriceTrip(t *testing.T) {
	// Create env
	env := test.NewTestEnv(t, container, "TestGetPriceTrip")
	assert.NotNil(t, env)

	client := env.Client
	db := env.Database

	// Mock mysql data
	mockUsers := mockUsers(client, 2)

	tests := []struct {
		name string
		info *TestGetPriceTripInfo
	}{
		{
			name: "[GetPriceTrip][Success] Return 200 - Return price of trip",
			info: &TestGetPriceTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				queryParams: map[string]string{
					"distance": "1",
				},
				output: &GetPriceTripResponse{
					Code: http.StatusOK,
					Data: &TypeResponse{
						Motor: 20000,
						Car:   30000,
					},
				},
			},
		},
		{
			name: "[GetPriceTrip][Fail] Return 400 - missing distance",
			info: &TestGetPriceTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				queryParams: map[string]string{},
				output: &GetPriceTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[GetPriceTrip][Fail] Return 400 - distance is not number",
			info: &TestGetPriceTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				queryParams: map[string]string{},
				output: &GetPriceTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doGetPriceTrip(t, tc.info, db)
		})
	}
}

func doGetPriceTrip(t *testing.T, info *TestGetPriceTripInfo, db *sql.DB) {
	req := test.BuildGetQuery(TripEndpoint, info.queryParams)
	rec := httptest.NewRecorder()
	c := NewTestEchoContext().NewContext(req, rec)
	if info.authInfo != nil {
		auth.SetAuthInfo(c, info.authInfo.UserID, info.authInfo.DriverID)
	}

	tripSrv := NewTripServer(db)
	err := tripSrv.GetPriceTrip(c)

	if info.output.Code == http.StatusOK {
		assert.Nil(t, err)
		var actual *GetPriceTripResponse
		err := json.NewDecoder(rec.Body).Decode(&actual)
		assert.Nil(t, err)
		assert.Equal(t, info.output.Code, actual.Code)
		assert.Equal(t, info.output.Data, actual.Data)
	} else {
		assert.Equal(t, info.output.Code, GetErrorCode(err))
	}
}

type TestListTripsInfo struct {
	authInfo    *auth.AuthInfo
	queryParams map[string]string
	output      *ListTripResponse
}

func TestListTrips(t *testing.T) {
	// Create env
	env := test.NewTestEnv(t, container, "TestListTrips")
	assert.NotNil(t, env)

	client := env.Client
	db := env.Database

	// Mock mysql data
	mockUsers := mockUsers(client, 2)
	mockDrivers := mockDrivers(client, 1)
	mockTrips := mockTrips(client, mockUsers[0].ID, mockDrivers[0].ID, 3)
	respData, _ := trip.DecodeManyTrip(mockTrips)

	tests := []struct {
		name string
		info *TestListTripsInfo
	}{
		{
			name: "[ListTrips][Success] Return 200 - Return all trips",
			info: &TestListTripsInfo{
				queryParams: map[string]string{},
				output: &ListTripResponse{
					Code: http.StatusOK,
					Data: respData,
				},
			},
		},
		{
			name: "[ListTrips][Success] Return 200 - Return trips of user[0]",
			info: &TestListTripsInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				queryParams: map[string]string{},
				output: &ListTripResponse{
					Code: http.StatusOK,
					Data: respData,
				},
			},
		},
		{
			name: "[ListTrips][Success] Return 200 - Return trips of user[1]",
			info: &TestListTripsInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[1].ID,
				},
				queryParams: map[string]string{},
				output: &ListTripResponse{
					Code: http.StatusOK,
					Data: []*trip.Trip{},
				},
			},
		},
		{
			name: "[ListTrips][Success] Return 200 - Return trips waiting",
			info: &TestListTripsInfo{
				authInfo: &auth.AuthInfo{
					DriverID: mockDrivers[0].ID,
				},
				queryParams: map[string]string{},
				output: &ListTripResponse{
					Code: http.StatusOK,
					Data: respData,
				},
			},
		},
		{
			name: "[ListTrips][Success] Return 200 - Return trips of user[0] where id=trip[0]",
			info: &TestListTripsInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				queryParams: map[string]string{
					"tripId": fmt.Sprint(mockTrips[0].ID),
				},
				output: &ListTripResponse{
					Code: http.StatusOK,
					Data: respData[:1],
				},
			},
		},
		{
			name: "[ListTrips][Success] Return 200 - Return trips of user[0] where status='waiting'",
			info: &TestListTripsInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				queryParams: map[string]string{
					"status": "waiting",
				},
				output: &ListTripResponse{
					Code: http.StatusOK,
					Data: respData,
				},
			},
		},
		{
			name: "[ListTrips][Success] Return 200 - Return trips of user[0] where rate=2",
			info: &TestListTripsInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				queryParams: map[string]string{
					"rate": "2",
				},
				output: &ListTripResponse{
					Code: http.StatusOK,
					Data: respData[1:2],
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doListTrips(t, tc.info, db)
		})
	}
}

func doListTrips(t *testing.T, info *TestListTripsInfo, db *sql.DB) {
	req := test.BuildGetQuery(TripEndpoint, info.queryParams)
	rec := httptest.NewRecorder()
	c := NewTestEchoContext().NewContext(req, rec)
	if info.authInfo != nil {
		auth.SetAuthInfo(c, info.authInfo.UserID, info.authInfo.DriverID)
	}

	tripSrv := NewTripServer(db)
	err := tripSrv.ListTrip(c)

	if info.output.Code == http.StatusOK {
		assert.Nil(t, err)
		var actual *ListTripResponse
		err := json.NewDecoder(rec.Body).Decode(&actual)
		assert.Nil(t, err)
		assert.Equal(t, info.output.Code, actual.Code)
		assert.Equal(t, len(info.output.Data), len(actual.Data))
	} else {
		assert.Equal(t, info.output.Code, GetErrorCode(err))
	}
}

type TestCreateTripInfo struct {
	authInfo *auth.AuthInfo
	body     map[string]interface{}
	output   *CreateTripResponse
}

func TestCreateTrip(t *testing.T) {
	// Create env
	env := test.NewTestEnv(t, container, "TestCreateTrip")
	assert.NotNil(t, env)

	client := env.Client
	db := env.Database

	// Mock mysql data
	mockUsers := mockUsers(client, 1)

	tests := []struct {
		name string
		info *TestCreateTripInfo
	}{
		{
			name: "[CreateTrip][Success] Return 200 - Return created trip with motor",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startX":        1,
					"startY":        1,
					"startLocation": "di",
					"endX":          1,
					"endY":          1,
					"endLocation":   "den",
					"type":          "motor",
					"distance":      1.2,
				},
				output: &CreateTripResponse{
					Code: http.StatusOK,
					Data: &trip.Trip{
						Status: "waiting",
						Price:  24000,
					},
				},
			},
		},
		{
			name: "[CreateTrip][Success] Return 200 - Return created trip with car",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startX":        1,
					"startY":        1,
					"startLocation": "di",
					"endX":          1,
					"endY":          1,
					"endLocation":   "den",
					"type":          "car",
					"distance":      1.2,
				},
				output: &CreateTripResponse{
					Code: http.StatusOK,
					Data: &trip.Trip{
						Status: "waiting",
						Price:  36000,
					},
				},
			},
		},
		{
			name: "[CreateTrip][Fail] Return 400 - Error startX not a number",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startX":        "test",
					"startY":        1,
					"startLocation": "di",
					"endX":          1,
					"endY":          1,
					"endLocation":   "den",
					"type":          "motor",
					"distance":      1.2,
				},
				output: &CreateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[CreateTrip][Fail] Return 400 - Error startY not a number",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startX":        1,
					"startY":        "test",
					"startLocation": "di",
					"endX":          1,
					"endY":          1,
					"endLocation":   "den",
					"type":          "motor",
					"distance":      1.2,
				},
				output: &CreateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[CreateTrip][Fail] Return 400 - Error endX not a number",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startX":        1,
					"startY":        1,
					"startLocation": "di",
					"endX":          "test",
					"endY":          1,
					"endLocation":   "den",
					"type":          "motor",
					"distance":      1.2,
				},
				output: &CreateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[CreateTrip][Fail] Return 400 - Error endY not a number",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startX":        1,
					"startY":        1,
					"startLocation": "di",
					"endX":          1,
					"endY":          "test",
					"endLocation":   "den",
					"type":          "motor",
					"distance":      1.2,
				},
				output: &CreateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[CreateTrip][Fail] Return 400 - Error price not a number",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startX":        1,
					"startY":        1,
					"startLocation": "di",
					"endX":          1,
					"endY":          1,
					"endLocation":   "den",
					"type":          "motor",
					"distance":      "test",
				},
				output: &CreateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[CreateTrip][Fail] Return 400 - invalid type",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startX":        1,
					"startY":        1,
					"startLocation": "di",
					"endX":          1,
					"endY":          1,
					"endLocation":   "den",
					"type":          "test",
					"distance":      1.2,
				},
				output: &CreateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[CreateTrip][Fail] Return 400 - missing require field type",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startX":        1,
					"startY":        1,
					"startLocation": "di",
					"endX":          1,
					"endY":          1,
					"endLocation":   "den",
					"distance":      1.2,
				},
				output: &CreateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[CreateTrip][Fail] Return 400 - missing require field startX",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startY":        1,
					"startLocation": "di",
					"endX":          1,
					"endY":          1,
					"endLocation":   "den",
					"type":          "motor",
					"distance":      1.2,
				},
				output: &CreateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[CreateTrip][Fail] Return 400 - missing require field startY",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startX":        1,
					"startLocation": "di",
					"endX":          1,
					"endY":          1,
					"endLocation":   "den",
					"type":          "motor",
					"distance":      1.2,
				},
				output: &CreateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[CreateTrip][Fail] Return 400 - missing require field startLocation",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startX":      1,
					"startY":      1,
					"endX":        1,
					"endY":        1,
					"endLocation": "den",
					"type":        "motor",
					"distance":    1.2,
				},
				output: &CreateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[CreateTrip][Fail] Return 400 - missing require field endX",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startX":        1,
					"startY":        1,
					"startLocation": "di",
					"endY":          1,
					"endLocation":   "den",
					"type":          "motor",
					"distance":      1.2,
				},
				output: &CreateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[CreateTrip][Fail] Return 400 - missing require field endY",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startX":        1,
					"startY":        1,
					"startLocation": "di",
					"endX":          1,
					"endLocation":   "den",
					"type":          "motor",
					"distance":      1.2,
				},
				output: &CreateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[CreateTrip][Fail] Return 400 - missing require field endLocation",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startX":        1,
					"startY":        1,
					"startLocation": "di",
					"endX":          1,
					"endY":          1,
					"type":          "motor",
					"distance":      1.2,
				},
				output: &CreateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[CreateTrip][Fail] Return 400 - missing require field distance",
			info: &TestCreateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"startX":        1,
					"startY":        1,
					"startLocation": "di",
					"endX":          1,
					"endY":          1,
					"endLocation":   "den",
					"type":          "motor",
				},
				output: &CreateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[CreateTrip][Fail] Return 400 - Error require login",
			info: &TestCreateTripInfo{
				body: map[string]interface{}{
					"startX":        1,
					"startY":        1,
					"startLocation": "di",
					"endX":          1,
					"endY":          1,
					"endLocation":   "den",
					"type":          "motor",
					"distance":      1,
				},
				output: &CreateTripResponse{
					Code: http.StatusUnauthorized,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doCreateTrip(t, tc.info, db)
		})
	}
}

func doCreateTrip(t *testing.T, info *TestCreateTripInfo, db *sql.DB) {
	req := test.BuildPostQuery(TripEndpoint, info.body)
	rec := httptest.NewRecorder()
	c := NewTestEchoContext().NewContext(req, rec)
	if info.authInfo != nil {
		auth.SetAuthInfo(c, info.authInfo.UserID, info.authInfo.DriverID)
	}

	tripSrv := NewTripServer(db)
	err := tripSrv.CreateTrip(c)

	if info.output.Code == http.StatusOK {
		assert.Nil(t, err)
		var actual *CreateTripResponse
		err := json.NewDecoder(rec.Body).Decode(&actual)
		assert.Nil(t, err)
		assert.Equal(t, info.output.Code, actual.Code)
		assert.Equal(t, info.output.Data.Status, actual.Data.Status)
		assert.Equal(t, info.output.Data.Price, actual.Data.Price)
	} else {
		assert.Equal(t, info.output.Code, GetErrorCode(err))
	}
}

type TestUpdateStatusTripInfo struct {
	authInfo  *auth.AuthInfo
	urlParams map[string]int
	body      map[string]interface{}
	output    *UpdateStatusTripResponse
}

func TestUpdateStatusTrip(t *testing.T) {
	// Create env
	env := test.NewTestEnv(t, container, "TestUpdateStatusTrip")
	assert.NotNil(t, env)

	client := env.Client
	db := env.Database

	// Mock mysql data
	mockUsers := mockUsers(client, 2)
	mockDrivers := mockDrivers(client, 2)
	mockTrips := mockTrips(client, mockUsers[0].ID, mockDrivers[0].ID, 3)

	tests := []struct {
		name string
		info *TestUpdateStatusTripInfo
	}{
		{
			name: "[UpdateStatusTrip][Success] Return 200 - user update status cancel success",
			info: &TestUpdateStatusTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[0].ID,
				},
				body: map[string]interface{}{
					"status": "cancel",
				},
				output: &UpdateStatusTripResponse{
					Code: http.StatusOK,
					Data: &trip.Trip{
						Status: "cancel",
					},
				},
			},
		},
		{
			name: "[UpdateStatusTrip][Fail] Return 400 - trip canceled",
			info: &TestUpdateStatusTripInfo{
				authInfo: &auth.AuthInfo{
					DriverID: mockDrivers[0].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[0].ID,
				},
				body: map[string]interface{}{
					"status": "done",
				},
				output: &UpdateStatusTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[UpdateStatusTrip][Success] Return 200 - driver update status accept success",
			info: &TestUpdateStatusTripInfo{
				authInfo: &auth.AuthInfo{
					DriverID: mockDrivers[0].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[1].ID,
				},
				body: map[string]interface{}{
					"status": "accept",
				},
				output: &UpdateStatusTripResponse{
					Code: http.StatusOK,
					Data: &trip.Trip{
						DriveID: mockDrivers[0].ID,
						Status:  "accept",
					},
				},
			},
		},
		{
			name: "[UpdateStatusTrip][Success] Return 200 - driver update status done success",
			info: &TestUpdateStatusTripInfo{
				authInfo: &auth.AuthInfo{
					DriverID: mockDrivers[0].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[1].ID,
				},
				body: map[string]interface{}{
					"status": "done",
				},
				output: &UpdateStatusTripResponse{
					Code: http.StatusOK,
					Data: &trip.Trip{
						DriveID: mockDrivers[0].ID,
						Status:  "done",
					},
				},
			},
		},
		{
			name: "[UpdateStatusTrip][Fail] Return 400 - trip done",
			info: &TestUpdateStatusTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[1].ID,
				},
				body: map[string]interface{}{
					"status": "cancel",
				},
				output: &UpdateStatusTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[UpdateStatusTrip][Fail] Return 400 - missing require field",
			info: &TestUpdateStatusTripInfo{
				authInfo: &auth.AuthInfo{
					DriverID: mockDrivers[0].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[2].ID,
				},
				body: map[string]interface{}{},
				output: &UpdateStatusTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[UpdateStatusTrip][Fail] Return 403 - driver can not update he does not have",
			info: &TestUpdateStatusTripInfo{
				authInfo: &auth.AuthInfo{
					DriverID: mockDrivers[1].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[2].ID,
				},
				body: map[string]interface{}{
					"status": "done",
				},
				output: &UpdateStatusTripResponse{
					Code: http.StatusForbidden,
				},
			},
		},
		{
			name: "[UpdateStatusTrip][Fail] Return 403 - user can not update he does not have",
			info: &TestUpdateStatusTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[1].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[2].ID,
				},
				body: map[string]interface{}{
					"status": "cancel",
				},
				output: &UpdateStatusTripResponse{
					Code: http.StatusForbidden,
				},
			},
		},
		{
			name: "[UpdateStatusTrip][Fail] Return 400 - user can not update done",
			info: &TestUpdateStatusTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[2].ID,
				},
				body: map[string]interface{}{
					"status": "done",
				},
				output: &UpdateStatusTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[UpdateStatusTrip][Fail] Return 400 - user can not update accept",
			info: &TestUpdateStatusTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[2].ID,
				},
				body: map[string]interface{}{
					"status": "accept",
				},
				output: &UpdateStatusTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[UpdateStatusTrip][Fail] Return 400 - user can not update strange status",
			info: &TestUpdateStatusTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[2].ID,
				},
				body: map[string]interface{}{
					"status": "test",
				},
				output: &UpdateStatusTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[UpdateStatusTrip][Fail] Return 400 - driver can not update waiting",
			info: &TestUpdateStatusTripInfo{
				authInfo: &auth.AuthInfo{
					DriverID: mockDrivers[0].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[2].ID,
				},
				body: map[string]interface{}{
					"status": "waiting",
				},
				output: &UpdateStatusTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[UpdateStatusTrip][Fail] Return 400 - driver can not update cancel",
			info: &TestUpdateStatusTripInfo{
				authInfo: &auth.AuthInfo{
					DriverID: mockDrivers[0].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[2].ID,
				},
				body: map[string]interface{}{
					"status": "cancel",
				},
				output: &UpdateStatusTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[UpdateStatusTrip][Fail] Return 400 - driver can not update strange status",
			info: &TestUpdateStatusTripInfo{
				authInfo: &auth.AuthInfo{
					DriverID: mockDrivers[0].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[2].ID,
				},
				body: map[string]interface{}{
					"status": "test",
				},
				output: &UpdateStatusTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[UpdateStatusTrip][Fail] Return 400 - need login before change status",
			info: &TestUpdateStatusTripInfo{
				urlParams: map[string]int{
					"tripId": mockTrips[0].ID,
				},
				body: map[string]interface{}{
					"status": "waiting",
				},
				output: &UpdateStatusTripResponse{
					Code: http.StatusUnauthorized,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doUpdateStatusTrip(t, tc.info, db)
		})
	}
}

func doUpdateStatusTrip(t *testing.T, info *TestUpdateStatusTripInfo, db *sql.DB) {
	req := test.BuildPatchRequest(TripEndpoint+"/:tripId", info.body)
	rec := httptest.NewRecorder()
	c := NewTestEchoContext().NewContext(req, rec)
	if info.authInfo != nil {
		auth.SetAuthInfo(c, info.authInfo.UserID, info.authInfo.DriverID)
	}
	if v, ok := info.urlParams["tripId"]; ok {
		c.SetParamNames("tripId")
		c.SetParamValues(fmt.Sprint(v))
	}

	tripSrv := NewTripServer(db)
	err := tripSrv.UpdateStatusTrip(c)

	if info.output.Code == http.StatusOK {
		assert.Nil(t, err)
		var actual *UpdateStatusTripResponse
		err := json.NewDecoder(rec.Body).Decode(&actual)
		assert.Nil(t, err)
		assert.Equal(t, info.output.Code, actual.Code)
		assert.Equal(t, info.output.Data.Status, actual.Data.Status)
		if info.output.Data.DriveID != 0 {
			assert.Equal(t, info.output.Data.DriveID, actual.Data.DriveID)
		}
	} else {
		assert.Equal(t, info.output.Code, GetErrorCode(err))
	}
}

type TestRateTripInfo struct {
	authInfo  *auth.AuthInfo
	urlParams map[string]int
	body      map[string]interface{}
	output    *RateTripResponse
}

func TestRateTrip(t *testing.T) {
	// Create env
	env := test.NewTestEnv(t, container, "TestRateTrip")
	assert.NotNil(t, env)

	client := env.Client
	db := env.Database

	// Mock mysql data
	mockUsers := mockUsers(client, 2)
	mockTrips := mockTrips(client, mockUsers[0].ID, 0, 1)

	tests := []struct {
		name string
		info *TestRateTripInfo
	}{
		{
			name: "[RateTrip][Success] Return 200 - user rate success",
			info: &TestRateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[0].ID,
				},
				body: map[string]interface{}{
					"rate": 1,
				},
				output: &RateTripResponse{
					Code: http.StatusOK,
					Data: &trip.Trip{
						Rate: 1,
					},
				},
			},
		},
		{
			name: "[RateTrip][Fail] Return 403 - user can not update user does not have",
			info: &TestRateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[1].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[0].ID,
				},
				body: map[string]interface{}{
					"rate": 1,
				},
				output: &RateTripResponse{
					Code: http.StatusForbidden,
				},
			},
		},
		{
			name: "[RateTrip][Fail] Return 401 - need login before",
			info: &TestRateTripInfo{
				urlParams: map[string]int{
					"tripId": mockTrips[0].ID,
				},
				body: map[string]interface{}{
					"rate": 1,
				},
				output: &RateTripResponse{
					Code: http.StatusUnauthorized,
				},
			},
		},
		{
			name: "[RateTrip][Fail] Return 400 - rate not a number",
			info: &TestRateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[1].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[0].ID,
				},
				body: map[string]interface{}{
					"rate": "test",
				},
				output: &RateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[RateTrip][Fail] Return 400 - missing require field",
			info: &TestRateTripInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[1].ID,
				},
				urlParams: map[string]int{
					"tripId": mockTrips[0].ID,
				},
				body: map[string]interface{}{},
				output: &RateTripResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doRateTrip(t, tc.info, db)
		})
	}
}

func doRateTrip(t *testing.T, info *TestRateTripInfo, db *sql.DB) {
	req := test.BuildPatchRequest(TripEndpoint+"/:tripId", info.body)
	rec := httptest.NewRecorder()
	c := NewTestEchoContext().NewContext(req, rec)
	if info.authInfo != nil {
		auth.SetAuthInfo(c, info.authInfo.UserID, info.authInfo.DriverID)
	}
	if v, ok := info.urlParams["tripId"]; ok {
		c.SetParamNames("tripId")
		c.SetParamValues(fmt.Sprint(v))
	}

	tripSrv := NewTripServer(db)
	err := tripSrv.RateTrip(c)

	if info.output.Code == http.StatusOK {
		assert.Nil(t, err)
		var actual *RateTripResponse
		err := json.NewDecoder(rec.Body).Decode(&actual)
		assert.Nil(t, err)
		assert.Equal(t, info.output.Code, actual.Code)
		assert.Equal(t, info.output.Data.Rate, actual.Data.Rate)
	} else {
		assert.Equal(t, info.output.Code, GetErrorCode(err))
	}
}

func mockTrips(client *ent.Client, userID, driverID, len int) []*ent.Trip {
	var ret []*ent.Trip
	for idx := 1; idx <= len; idx++ {
		trip := &ent.Trip{
			UserID:        userID,
			DriverID:      driverID,
			StartX:        float64(idx),
			StartY:        float64(idx),
			StartLocation: "di",
			EndX:          float64(idx),
			EndY:          float64(idx),
			EndLocation:   "den",
			Type:          entTrip.TypeMotor,
			Price:         20,
			Distance:      float64(idx),
			Rate:          idx % 5,
		}
		ret = append(ret, trip)
	}

	bulk := make([]*ent.TripCreate, 0)
	for _, trip := range ret {
		q := client.Trip.Create().
			SetUserID(trip.UserID).
			SetStartX(trip.StartX).
			SetStartY(trip.StartY).
			SetStartLocation(trip.StartLocation).
			SetEndX(trip.EndX).
			SetEndY(trip.EndY).
			SetEndLocation(trip.EndLocation).
			SetPrice(trip.Price).
			SetType(trip.Type).
			SetDistance(trip.Distance).
			SetRate(trip.Rate)
		if trip.DriverID != 0 {
			q.SetDriverID(trip.DriverID)
		}
		bulk = append(bulk, q)
	}

	mockData, _ := client.Trip.CreateBulk(bulk...).
		Save(context.Background())
	return mockData
}
