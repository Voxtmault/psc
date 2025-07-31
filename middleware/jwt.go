package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

const (
	JWTClaimsUserID    = "user_id"
	JWTClaimsCompanyID = "company_id"
)

func VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		bearerToken := c.Request().Header.Get("Authorization")
		if bearerToken == "" {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
				"status_code": http.StatusBadRequest,
				"message":     "Bad Request",
				"error":       "Missing Authorization Header",
			})
		}

		if !strings.HasPrefix(bearerToken, "Bearer ") {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
				"status_code": http.StatusBadRequest,
				"message":     "Bad Request",
				"error":       "Invalid Field Format Authorization",
			})
		}

		token, err := jwt.Parse(strings.TrimPrefix(bearerToken, "Bearer "), func(token *jwt.Token) (interface{}, error) {
			return []byte("jwt_secret"), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			if errors.Is(err, jwt.ErrTokenMalformed) {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"status_code": http.StatusUnauthorized,
					"message":     "Unauthorized",
					"error":       "Invalid JWT Format",
				})
			} else if errors.Is(err, jwt.ErrTokenExpired) {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"status_code": http.StatusUnauthorized,
					"message":     "Unauthorized",
					"error":       "Expired Token",
				})
			} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"status_code": http.StatusUnauthorized,
					"message":     "Unauthorized",
					"error":       "Invalid JWT Signature",
				})
			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, map[string]interface{}{
					"status_code": http.StatusInternalServerError,
					"message":     "General Error",
					"error":       eris.Wrap(err, "parsing jwt token").Error(),
				})
			}
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]interface{}{
				"status_code": http.StatusInternalServerError,
				"message":     "General Error",
				"error":       eris.Wrap(err, "failed to parse claims").Error(),
			})
		}

		c.Set(JWTClaimsUserID, claims["user_id"].(string))
		c.Set(JWTClaimsCompanyID, claims["company_id"].(string))

		return next(c)
	}
}
