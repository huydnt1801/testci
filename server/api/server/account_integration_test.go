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
	"testing"

	"github.com/huydnt1801/chuyende/api/server/middleware/auth"
	"github.com/huydnt1801/chuyende/internal/driver"
	"github.com/huydnt1801/chuyende/internal/ent"
	"github.com/huydnt1801/chuyende/internal/ent/vehicledriver"
	"github.com/huydnt1801/chuyende/internal/user"
	"github.com/huydnt1801/chuyende/test"
	"github.com/stretchr/testify/assert"
)

const (
	AccountEndpoint string = "/api/v1/account"
)

func mockUsers(client *ent.Client, len int) []*ent.User {
	var ret []*ent.User
	passwordHasher := &user.BcryptPasswordHasher{}
	for idx := 1; idx <= len; idx++ {
		user := &ent.User{
			PhoneNumber: "012345678" + fmt.Sprint(idx),
			Confirmed:   true,
			FullName:    "test user " + fmt.Sprint(idx),
			Password:    "123456",
		}
		ret = append(ret, user)
	}
	bulk := make([]*ent.UserCreate, 0)
	for _, user := range ret {
		hashPass, _ := passwordHasher.HashPassword(user.Password)
		q := client.User.Create().
			SetFullName(user.FullName).
			SetConfirmed(user.Confirmed).
			SetPhoneNumber(user.PhoneNumber).
			SetPassword(hashPass)
		bulk = append(bulk, q)
	}

	mockData, _ := client.User.CreateBulk(bulk...).
		Save(context.Background())
	return mockData
}

func mockDrivers(client *ent.Client, len int) []*ent.VehicleDriver {
	var ret []*ent.VehicleDriver
	for idx := 1; idx <= len; idx++ {
		user := &ent.VehicleDriver{
			PhoneNumber: "098765432" + fmt.Sprint(idx),
			FullName:    "test driver " + fmt.Sprint(idx),
			Password:    "123456",
			License:     vehicledriver.LicenseMotor,
		}
		ret = append(ret, user)
	}

	bulk := make([]*ent.VehicleDriverCreate, 0)
	for _, user := range ret {
		q := client.VehicleDriver.Create().
			SetFullName(user.FullName).
			SetPhoneNumber(user.PhoneNumber).
			SetPassword(user.Password).
			SetLicense(user.License)
		bulk = append(bulk, q)
	}

	mockData, _ := client.VehicleDriver.CreateBulk(bulk...).
		Save(context.Background())
	return mockData
}

type TestLoginInfo struct {
	body   map[string]interface{}
	output *LoginResponse
}

func TestLogin(t *testing.T) {
	// Create env
	env := test.NewTestEnv(t, container, "TestLogin")
	assert.NotNil(t, env)

	client := env.Client
	db := env.Database

	// Mock mysql data
	mockUsers := mockUsers(client, 1)

	tests := []struct {
		name string
		info *TestLoginInfo
	}{
		{
			name: "[TestLogin][Success] Return 200 - Return login success",
			info: &TestLoginInfo{
				body: map[string]interface{}{
					"phoneNumber": mockUsers[0].PhoneNumber,
					"password":    "123456",
				},
				output: &LoginResponse{
					Code: http.StatusOK,
					Data: &user.User{
						ID: mockUsers[0].ID,
					},
				},
			},
		},
		{
			name: "[TestLogin][Failed] Return 400 - Invalid password",
			info: &TestLoginInfo{
				body: map[string]interface{}{
					"phoneNumber": mockUsers[0].PhoneNumber,
					"password":    "1234566",
				},
				output: &LoginResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[TestLogin][Failed] Return 400 - Invalid phone number",
			info: &TestLoginInfo{
				body: map[string]interface{}{
					"phoneNumber": mockUsers[0].PhoneNumber + "1",
					"password":    "123456",
				},
				output: &LoginResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[TestLogin][Failed] Return 400 - missing phone number",
			info: &TestLoginInfo{
				body: map[string]interface{}{
					"password": "123456",
				},
				output: &LoginResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[TestLogin][Failed] Return 400 - missing password",
			info: &TestLoginInfo{
				body: map[string]interface{}{
					"phoneNumber": mockUsers[0].PhoneNumber,
				},
				output: &LoginResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doLogin(t, tc.info, db)
		})
	}
}

func doLogin(t *testing.T, info *TestLoginInfo, db *sql.DB) {
	req := test.BuildPostQuery(AccountEndpoint+"/login", info.body)
	rec := httptest.NewRecorder()
	c := NewTestEchoContext().NewContext(req, rec)

	accountSv := NewAccountServer(db)
	err := accountSv.Login(c)

	if info.output.Code == http.StatusOK {
		assert.Nil(t, err)
		var actual *LoginResponse
		err := json.NewDecoder(rec.Body).Decode(&actual)
		assert.Nil(t, err)
		assert.Equal(t, info.output.Code, actual.Code)
		assert.Equal(t, info.output.Data.ID, actual.Data.ID)
	} else {
		assert.Equal(t, info.output.Code, GetErrorCode(err))
	}
}

type TestLoginDriverInfo struct {
	body   map[string]interface{}
	output *LoginDriverResponse
}

func TestLoginDriver(t *testing.T) {
	// Create env
	env := test.NewTestEnv(t, container, "TestLoginDriver")
	assert.NotNil(t, env)

	client := env.Client
	db := env.Database

	// Mock mysql data
	mockDrivers := mockDrivers(client, 1)

	tests := []struct {
		name string
		info *TestLoginDriverInfo
	}{
		{
			name: "[TestLoginDriver][Success] Return 200 - Return login success",
			info: &TestLoginDriverInfo{
				body: map[string]interface{}{
					"phoneNumber": mockDrivers[0].PhoneNumber,
					"password":    "123456",
				},
				output: &LoginDriverResponse{
					Code: http.StatusOK,
					Data: &driver.Driver{
						ID: mockDrivers[0].ID,
					},
				},
			},
		},
		{
			name: "[TestLogin][Failed] Return 400 - Invalid password",
			info: &TestLoginDriverInfo{
				body: map[string]interface{}{
					"phoneNumber": mockDrivers[0].PhoneNumber,
					"password":    "1234561",
				},
				output: &LoginDriverResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[TestLogin][Failed] Return 400 - Invalid phone number",
			info: &TestLoginDriverInfo{
				body: map[string]interface{}{
					"phoneNumber": mockDrivers[0].PhoneNumber + "1",
					"password":    "123456",
				},
				output: &LoginDriverResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[TestLogin][Failed] Return 400 - missing phone number",
			info: &TestLoginDriverInfo{
				body: map[string]interface{}{
					"password": "123456",
				},
				output: &LoginDriverResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[TestLogin][Failed] Return 400 - missing phone number",
			info: &TestLoginDriverInfo{
				body: map[string]interface{}{
					"phoneNumber": mockDrivers[0].PhoneNumber,
				},
				output: &LoginDriverResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doLoginDriver(t, tc.info, db)
		})
	}
}

func doLoginDriver(t *testing.T, info *TestLoginDriverInfo, db *sql.DB) {
	req := test.BuildPostQuery(AccountEndpoint+"/login/driver", info.body)
	rec := httptest.NewRecorder()
	c := NewTestEchoContext().NewContext(req, rec)

	accountSv := NewAccountServer(db)
	err := accountSv.LoginDriver(c)

	if info.output.Code == http.StatusOK {
		assert.Nil(t, err)
		var actual *LoginDriverResponse
		err := json.NewDecoder(rec.Body).Decode(&actual)
		assert.Nil(t, err)
		assert.Equal(t, info.output.Code, actual.Code)
		assert.Equal(t, info.output.Data.ID, actual.Data.ID)
	} else {
		assert.Equal(t, info.output.Code, GetErrorCode(err))
	}
}

type TestRegisterInfo struct {
	body   map[string]interface{}
	output *RegisterResponse
}

func TestRegister(t *testing.T) {
	// Create env
	env := test.NewTestEnv(t, container, "TestRegister")
	assert.NotNil(t, env)

	client := env.Client
	db := env.Database

	// Mock mysql data
	mockUsers := mockUsers(client, 1)

	tests := []struct {
		name string
		info *TestRegisterInfo
	}{
		{
			name: "[TestRegister][Success] Return 200 - Return register success",
			info: &TestRegisterInfo{
				body: map[string]interface{}{
					"phoneNumber": "0987654321",
					"fullName":    "test register",
					"password":    "123456",
					"password2":   "123456",
				},
				output: &RegisterResponse{
					Code: http.StatusOK,
				},
			},
		},
		{
			name: "[TestRegister][Fail] Return 400 - missing phone number",
			info: &TestRegisterInfo{
				body: map[string]interface{}{
					"fullName":  "test register",
					"password":  "123456",
					"password2": "123456",
				},
				output: &RegisterResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[TestRegister][Fail] Return 400 - missing full name",
			info: &TestRegisterInfo{
				body: map[string]interface{}{
					"phoneNumber": "0987654321",
					"password":    "123456",
					"password2":   "123456",
				},
				output: &RegisterResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[TestRegister][Fail] Return 400 - missing password",
			info: &TestRegisterInfo{
				body: map[string]interface{}{
					"phoneNumber": "0987654321",
					"fullName":    "test register",
					"password2":   "123456",
				},
				output: &RegisterResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[TestRegister][Fail] Return 400 - missing password2",
			info: &TestRegisterInfo{
				body: map[string]interface{}{
					"phoneNumber": "0987654321",
					"fullName":    "test register",
					"password":    "123456",
				},
				output: &RegisterResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[TestRegister][Fail] Return 400 - invalid phone number",
			info: &TestRegisterInfo{
				body: map[string]interface{}{
					"phoneNumber": "09876543test",
					"fullName":    "test register",
					"password":    "123456",
					"password2":   "123456",
				},
				output: &RegisterResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[TestRegister][Fail] Return 400 - invalid password",
			info: &TestRegisterInfo{
				body: map[string]interface{}{
					"phoneNumber": "09876543test",
					"fullName":    "test register",
					"password":    "1234567",
					"password2":   "1234567",
				},
				output: &RegisterResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[TestRegister][Fail] Return 400 - user existed",
			info: &TestRegisterInfo{
				body: map[string]interface{}{
					"phoneNumber": mockUsers[0].PhoneNumber,
					"fullName":    "test register",
					"password":    "123456",
					"password2":   "123456",
				},
				output: &RegisterResponse{
					Code: http.StatusConflict,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doRegister(t, tc.info, db)
		})
	}
}

func doRegister(t *testing.T, info *TestRegisterInfo, db *sql.DB) {
	req := test.BuildPostQuery(AccountEndpoint+"/register", info.body)
	rec := httptest.NewRecorder()
	c := NewTestEchoContext().NewContext(req, rec)

	accountSv := NewAccountServer(db)
	err := accountSv.Register(c)

	if info.output.Code == http.StatusOK {
		assert.Nil(t, err)
		var actual *RegisterResponse
		err := json.NewDecoder(rec.Body).Decode(&actual)
		assert.Nil(t, err)
		assert.Equal(t, info.output.Code, actual.Code)
	} else {
		assert.Equal(t, info.output.Code, GetErrorCode(err))
	}
}

type TestCheckPhoneInfo struct {
	body   map[string]interface{}
	output *CheckPhoneResponse
}

func TestCheckPhone(t *testing.T) {
	// Create env
	env := test.NewTestEnv(t, container, "TestCheckPhone")
	assert.NotNil(t, env)

	client := env.Client
	db := env.Database

	// Mock mysql data
	mockUsers := mockUsers(client, 1)

	tests := []struct {
		name string
		info *TestCheckPhoneInfo
	}{
		{
			name: "[TestCheckPhone][Success] Return 200 - Return user if exist",
			info: &TestCheckPhoneInfo{
				body: map[string]interface{}{
					"phoneNumber": mockUsers[0].PhoneNumber,
				},
				output: &CheckPhoneResponse{
					Code: http.StatusOK,
					Data: &user.User{
						ID: mockUsers[0].ID,
					},
				},
			},
		},
		{
			name: "[TestCheckPhone][Success] Return 200 - Return empty user if not exist",
			info: &TestCheckPhoneInfo{
				body: map[string]interface{}{
					"phoneNumber": "0987654321",
				},
				output: &CheckPhoneResponse{
					Code: http.StatusOK,
					Data: nil,
				},
			},
		},
		{
			name: "[TestCheckPhone][Fail] Return 400 - missing phone",
			info: &TestCheckPhoneInfo{
				body: map[string]interface{}{},
				output: &CheckPhoneResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[TestCheckPhone][Fail] Return 400 - invalid phone",
			info: &TestCheckPhoneInfo{
				body: map[string]interface{}{
					"phoneNumber": "098765432test",
				},
				output: &CheckPhoneResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doCheckPhone(t, tc.info, db)
		})
	}
}

func doCheckPhone(t *testing.T, info *TestCheckPhoneInfo, db *sql.DB) {
	req := test.BuildPostQuery(AccountEndpoint+"/phone", info.body)
	rec := httptest.NewRecorder()
	c := NewTestEchoContext().NewContext(req, rec)

	accountSv := NewAccountServer(db)
	err := accountSv.CheckPhone(c)

	if info.output.Code == http.StatusOK {
		assert.Nil(t, err)
		var actual *CheckPhoneResponse
		err := json.NewDecoder(rec.Body).Decode(&actual)
		assert.Nil(t, err)
		assert.Equal(t, info.output.Code, actual.Code)
		if info.output.Data != nil {
			assert.Equal(t, info.output.Data.ID, actual.Data.ID)
		} else {
			assert.Nil(t, actual.Data)
		}
	} else {
		assert.Equal(t, info.output.Code, GetErrorCode(err))
	}
}

type TestUpdateInfoInfo struct {
	authInfo  *auth.AuthInfo
	body      map[string]interface{}
	urlParams map[string]int
	output    *UpdateInfoResponse
}

func TestUpdateInfo(t *testing.T) {
	// Create env
	env := test.NewTestEnv(t, container, "TestUpdateInfo")
	assert.NotNil(t, env)

	client := env.Client
	db := env.Database

	// Mock mysql data
	mockUsers := mockUsers(client, 1)

	tests := []struct {
		name string
		info *TestUpdateInfoInfo
	}{
		{
			name: "[TestUpdateInfo][Success] Return 200 - Return updated user",
			info: &TestUpdateInfoInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"password": "123123",
					"imageUrl": "test",
					"fullName": "test",
				},
				urlParams: map[string]int{
					"userId": mockUsers[0].ID,
				},
				output: &UpdateInfoResponse{
					Code: http.StatusOK,
					Data: &user.User{
						ImageURL: "test",
						FullName: "test",
					},
				},
			},
		},
		{
			name: "[TestUpdateInfo][Fail] Return 400 - invalid password",
			info: &TestUpdateInfoInfo{
				authInfo: &auth.AuthInfo{
					UserID: mockUsers[0].ID,
				},
				body: map[string]interface{}{
					"password": "1231213",
					"imageUrl": "test",
					"fullName": "test",
				},
				urlParams: map[string]int{
					"userId": mockUsers[0].ID,
				},
				output: &UpdateInfoResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
		{
			name: "[TestUpdateInfo][Fail] Return 403 - forbidden",
			info: &TestUpdateInfoInfo{
				body: map[string]interface{}{
					"password": "123123",
					"imageUrl": "test",
					"fullName": "test",
				},
				urlParams: map[string]int{
					"userId": mockUsers[0].ID,
				},
				output: &UpdateInfoResponse{
					Code: http.StatusForbidden,
				},
			},
		},
		{
			name: "[TestUpdateInfo][Fail] Return 400 - missing userID",
			info: &TestUpdateInfoInfo{
				body: map[string]interface{}{
					"password": "123123",
					"imageUrl": "test",
					"fullName": "test",
				},
				output: &UpdateInfoResponse{
					Code: http.StatusBadRequest,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doUpdateInfo(t, tc.info, db)
		})
	}
}

func doUpdateInfo(t *testing.T, info *TestUpdateInfoInfo, db *sql.DB) {
	req := test.BuildPatchRequest(AccountEndpoint+"/:userId", info.body)
	rec := httptest.NewRecorder()
	c := NewTestEchoContext().NewContext(req, rec)
	if info.authInfo != nil {
		auth.SetAuthInfo(c, info.authInfo.UserID, info.authInfo.DriverID)
	}
	if v, ok := info.urlParams["userId"]; ok {
		c.SetParamNames("userId")
		c.SetParamValues(fmt.Sprint(v))
	}

	accountSv := NewAccountServer(db)
	err := accountSv.UpdateInfo(c)

	if info.output.Code == http.StatusOK {
		assert.Nil(t, err)
		var actual *UpdateInfoResponse
		err := json.NewDecoder(rec.Body).Decode(&actual)
		assert.Nil(t, err)
		assert.Equal(t, info.output.Code, actual.Code)
		assert.Equal(t, info.output.Data.ImageURL, actual.Data.ImageURL)
		assert.Equal(t, info.output.Data.FullName, actual.Data.FullName)
	} else {
		assert.Equal(t, info.output.Code, GetErrorCode(err))
	}
}
