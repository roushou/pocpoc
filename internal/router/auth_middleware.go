package router

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/roushou/pocpoc/internal/security"
)

type userIDContextKey string

const userIDKey userIDContextKey = "userID"

// jwtCookieName is the cookie name in which JWTs are stored.
const jwtCookieName = "token"

func AuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			rc, ok := ctx.(*routerContext)
			if !ok {
				return echo.ErrInternalServerError
			}

			cookie, err := ctx.Cookie(jwtCookieName)
			if err != nil {
				return echo.ErrUnauthorized
			}

			claims := &JWTClaims{}
			_, err = security.ParseJWTWithClaims(cookie.Value, claims, secret)
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					return echo.ErrUnauthorized
				}
				return echo.ErrBadRequest
			}

			if claims.UserID == uuid.Nil {
				return echo.ErrUnauthorized
			}

			// Set user ID in context for use in handlers
			rc.Set(string(userIDKey), claims.UserID)
			return next(rc)
		}
	}
}
