package router

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/roushou/pocpoc/internal/models"
	"github.com/roushou/pocpoc/internal/security"
	"gorm.io/gorm"
)

// TODO: Load secret key and expiration time from config
const jwtSecretKey = "bigsecret"

// Setting it high for simplicity
const jwtExpiration = 24 * time.Hour

func bindAuthRouter(router *echo.Group) {
	group := router.Group("/auth")
	group.POST("/sign-in", signIn)
	group.POST("/sign-up", signUp)
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"userID"`
}

func signIn(ctx echo.Context) error {
	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := ctx.Bind(&payload); err != nil {
		return echo.ErrBadRequest
	}
	if len(payload.Username) == 0 || len(payload.Password) == 0 {
		return echo.ErrBadRequest
	}

	hashedPassword, err := security.HashPassword(payload.Password)
	if err != nil {
		return echo.ErrInternalServerError
	}

	db := ctx.(*routerContext).GetDatabase()

	user := &models.User{
		Username:     payload.Username,
		PasswordHash: hashedPassword,
	}
	tx := db.First(user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	token, err := security.NewJWT(JWTClaims{UserID: user.ID}, jwtSecretKey, jwtExpiration)
	if err != nil {
		return echo.ErrInternalServerError
	}

	ctx.SetCookie(&http.Cookie{
		Name:     jwtCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	return ctx.String(http.StatusOK, "OK")
}

func signUp(ctx echo.Context) error {
	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := ctx.Bind(&payload); err != nil {
		return echo.ErrBadRequest
	}

	hashedPassword, err := security.HashPassword(payload.Password)
	if err != nil {
		return echo.ErrInternalServerError
	}

	db := ctx.(*routerContext).GetDatabase()

	user := &models.User{
		Username:     payload.Username,
		PasswordHash: hashedPassword,
	}
	tx := db.Create(user)
	if tx.Error != nil {
		return echo.ErrInternalServerError
	}

	token, err := security.NewJWT(&JWTClaims{UserID: user.ID}, jwtSecretKey, jwtExpiration)
	if err != nil {
		return echo.ErrInternalServerError
	}

	ctx.SetCookie(&http.Cookie{
		Name:     jwtCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	return ctx.JSON(http.StatusCreated, map[string]string{"message": "success"})
}
