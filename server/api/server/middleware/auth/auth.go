package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type AuthInfo struct {
	UserID   int
	DriverID int
}

func Middleware(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()

			if strings.HasPrefix(r.URL.Path, "/api/v1/accounts/register") ||
				strings.HasPrefix(r.URL.Path, "/api/v1/accounts/login") ||
				strings.HasPrefix(r.URL.Path, "/api/v1/accounts/confirm") ||
				strings.HasPrefix(r.URL.Path, "/api/v1/accounts/resend") ||
				strings.HasPrefix(r.URL.Path, "/api/v1/accounts/phone") {
				return next(c)
			}
			sess := getSession(c)
			val, found := sess.Values["accessToken"]
			if found {
				token := fmt.Sprintf("%v", val)
				usrID, drivID, err := validateJWT(token)
				if err != nil {
					return nil
				}
				SetAuthInfo(c, usrID, drivID)
				return next(c)
			}
			return echo.NewHTTPError(http.StatusUnauthorized, "Yêu cầu đăng nhập")
		}
	}
}

func LoginUser(c echo.Context, usrID int, drivID int) error {
	sess := getSession(c)
	token, err := generateJWT(usrID, drivID)
	if err != nil {
		return err
	}
	sess.Values["accessToken"] = token

	saveSession(c)
	SetAuthInfo(c, usrID, drivID)
	return nil
}

func LogoutUser(c echo.Context) {
	sess := getSession(c)
	delete(sess.Values, "accessToken")
	sess.Options.MaxAge = -1
	saveSession(c)
}

func getSession(c echo.Context) *sessions.Session {
	sess, err := session.Get("auth", c)
	if err != nil {
		return &sessions.Session{
			Values:  map[interface{}]interface{}{},
			Options: &sessions.Options{},
		}
	}
	return sess
}

func saveSession(c echo.Context) {
	// skip if hook is already added
	if added, ok := c.Get(sessionSaveHookAdded).(bool); ok && added {
		return
	}
	c.Response().Before(func() {
		sess := getSession(c)
		sess.Options.HttpOnly = true
		sess.Options.Secure = true
		sess.Options.SameSite = http.SameSiteNoneMode
		// Support testing
		if sess.Name() == "auth" {
			sess.Save(c.Request(), c.Response())
		}
	})
	c.Set(sessionSaveHookAdded, true)

}

func SetAuthInfo(c echo.Context, usr int, driv int) {
	c.Set(authInfoKey, AuthInfo{
		UserID:   usr,
		DriverID: driv,
	})

}
func GetAuthInfo(c echo.Context) (AuthInfo, bool) {
	val := c.Get(authInfoKey)
	if val == nil {
		return AuthInfo{}, false
	}
	info, ok := val.(AuthInfo)
	if !ok {
		return AuthInfo{}, false
	}
	return info, true
}

type MyCustomClaims struct {
	UserID   int `json:"userID"`
	DriverID int `json:"driverID"`
	jwt.StandardClaims
}

func generateJWT(userID, driverID int) (string, error) {
	claims := MyCustomClaims{
		userID,
		driverID,
		jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("AllYourBase"))
	if err != nil {
		return "Signing Error", err
	}

	return tokenString, nil
}

func validateJWT(accessToken string) (int, int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims.UserID, claims.DriverID, nil
	} else {
		return 0, 0, err
	}
}

const (
	authInfoKey          string = "authInfo"
	sessionSaveHookAdded string = "authSessionSaveHookAdded"
)
