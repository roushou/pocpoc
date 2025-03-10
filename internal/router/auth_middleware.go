package router

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/roushou/pocpoc/internal/models"
	"github.com/roushou/pocpoc/internal/security"
)

type userIDContextKey string

const userIDKey userIDContextKey = "userID"

// jwtCookieName is the cookie name in which JWTs are stored.
const jwtCookieName = "token"

type authUser struct {
	UserID uuid.UUID
	Role   models.Role
}

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

			// Set auth user in context for use in handlers
			rc.Set(string(userIDKey), authUser{UserID: claims.UserID, Role: claims.Role})
			return next(rc)
		}
	}
}
