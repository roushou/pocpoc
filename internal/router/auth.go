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
	group.POST("/owners/sign-up", signUpOwner)
	group.POST("/owners/sign-in", signInOwner)
	group.POST("/staff/sign-in", signInStaff)
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID   `json:"userID"`
	Role   models.Role `json:"role"`
}

func signInOwner(ctx echo.Context) error {
	payload := struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}{}
	if err := ctx.Bind(&payload); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(&payload); err != nil {
		return echo.ErrBadRequest
	}

	hashedPassword, err := security.HashPassword(payload.Password)
	if err != nil {
		return echo.ErrInternalServerError
	}

	db := ctx.(*routerContext).GetDatabase()

	owner := &models.Owner{
		Username:     payload.Username,
		PasswordHash: hashedPassword,
	}
	tx := db.Connection.First(owner)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	claims := JWTClaims{UserID: owner.ID, Role: models.RoleOwner}
	token, err := security.NewJWT(claims, jwtSecretKey, jwtExpiration)
	if err != nil {
		return echo.ErrInternalServerError
	}

	ctx.SetCookie(&http.Cookie{
		Name:     jwtCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}

func signInStaff(ctx echo.Context) error {
	payload := struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}{}
	if err := ctx.Bind(&payload); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(&payload); err != nil {
		return echo.ErrBadRequest
	}

	hashedPassword, err := security.HashPassword(payload.Password)
	if err != nil {
		return echo.ErrInternalServerError
	}

	db := ctx.(*routerContext).GetDatabase()

	staff := &models.Staff{
		Username:     payload.Username,
		PasswordHash: hashedPassword,
	}
	tx := db.Connection.First(staff)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	claims := JWTClaims{UserID: staff.ID, Role: models.RoleStaff}
	token, err := security.NewJWT(claims, jwtSecretKey, jwtExpiration)
	if err != nil {
		return echo.ErrInternalServerError
	}

	ctx.SetCookie(&http.Cookie{
		Name:     jwtCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}

func signUpOwner(ctx echo.Context) error {
	payload := struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}{}
	if err := ctx.Bind(&payload); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(&payload); err != nil {
		return err
	}

	hashedPassword, err := security.HashPassword(payload.Password)
	if err != nil {
		return echo.ErrInternalServerError
	}

	db := ctx.(*routerContext).GetDatabase()

	user := &models.Owner{
		Username:     payload.Username,
		PasswordHash: hashedPassword,
	}
	tx := db.Connection.Create(user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
			return echo.ErrConflict
		}
		return echo.ErrInternalServerError
	}

	claims := &JWTClaims{UserID: user.ID, Role: models.RoleOwner}
	token, err := security.NewJWT(claims, jwtSecretKey, jwtExpiration)
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
