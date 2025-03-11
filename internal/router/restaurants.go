package router

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/roushou/pocpoc/internal/models"
	"github.com/roushou/pocpoc/internal/security"
	"gorm.io/gorm"
)

func bindRestaurantsRouter(router *echo.Group) {
	group := router.Group("/restaurants")
	group.GET("/:restaurant_id", getRestaurantById)
	group.POST("", registerRestaurant)
	group.POST("/:restaurant_id/staff", registerStaff)
}

func getRestaurantById(ctx echo.Context) error {
	_, err := getAuthUser(ctx)
	if err != nil {
		return echo.ErrUnauthorized
	}

	restaurantID, err := uuid.Parse(ctx.Param("restaurant_id"))
	if err != nil {
		return echo.ErrBadRequest
	}

	restaurant := &models.Restaurant{}

	db := ctx.(*routerContext).GetDatabase()
	// Need to do this way to prevent SQL injections
	tx := db.First(restaurant, "id = ?", restaurantID)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	return ctx.JSON(http.StatusOK, restaurant)
}

func registerRestaurant(ctx echo.Context) error {
	authUser, err := getAuthUser(ctx)
	if err != nil {
		return echo.ErrUnauthorized
	}

	payload := struct {
		Name string `json:"name"`
	}{}
	if err := ctx.Bind(&payload); err != nil {
		return echo.ErrBadRequest
	}
	if len(payload.Name) == 0 {
		return echo.ErrBadRequest
	}

	db := ctx.(*routerContext).GetDatabase()

	restaurant := &models.Restaurant{
		Name:    payload.Name,
		OwnerID: authUser.UserID,
	}
	tx := db.Create(restaurant)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
			return echo.ErrConflict
		}
		return echo.ErrInternalServerError
	}

	return ctx.JSON(http.StatusCreated, restaurant)
}

func registerStaff(ctx echo.Context) error {
	authUser, err := getAuthUser(ctx)
	if err != nil {
		return echo.ErrUnauthorized
	}
	// Only owners can register staff
	if authUser.Role != models.RoleOwner {
		return echo.ErrUnauthorized
	}

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

	restaurantID, err := uuid.Parse(ctx.Param("restaurant_id"))
	if err != nil {
		return echo.ErrBadRequest
	}

	db := ctx.(*routerContext).GetDatabase()

	// check if owner is registering staff for their own restaurant
	restaurant := &models.Restaurant{
		ID:      restaurantID,
		OwnerID: authUser.UserID,
	}
	if err := db.First(restaurant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrUnauthorized
		}
		return echo.ErrInternalServerError
	}

	hashedPassword, err := security.HashPassword(payload.Password)
	if err != nil {
		return echo.ErrInternalServerError
	}

	staff := &models.Staff{
		RestaurantID: restaurantID,
		Username:     payload.Username,
		PasswordHash: hashedPassword,
	}
	tx := db.Create(staff)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
			return echo.ErrConflict
		}
		return echo.ErrInternalServerError
	}

	claims := &JWTClaims{UserID: staff.ID, Role: models.RoleStaff}
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
