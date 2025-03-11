package router

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/roushou/pocpoc/internal/models"
	"gorm.io/gorm"
)

func bindProductsRouter(router *echo.Group) {
	group := router.Group("/restaurants/:restaurant_id/products")
	group.GET("", getProducts)
	group.POST("", registerProduct)
}

// ProductResponse maps fields of Product model we are willing to expose.
type ProductResponse struct {
	ProductID    uuid.UUID `json:"product_id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func getProducts(ctx echo.Context) error {
	authUser, err := getAuthUser(ctx)
	if err != nil {
		return echo.ErrUnauthorized
	}

	restaurantID, err := uuid.Parse(ctx.Param("restaurant_id"))
	if err != nil {
		return echo.ErrBadRequest
	}

	db := ctx.(*routerContext).GetDatabase()

	// Staff can only retrieve products of their own restaurant
	staff := &models.Staff{ID: authUser.UserID, RestaurantID: restaurantID}
	if err := db.Connection.First(staff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	rows := make([]models.Product, 0)
	if err := db.
		Connection.
		Select("id", "restaurant_id", "title", "description").
		Where("restaurant_id = ?", restaurantID).
		Find(&rows).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	products := make([]ProductResponse, 0, len(rows))

	for _, product := range rows {
		products = append(products, ProductResponse{
			ProductID:    product.ID,
			RestaurantID: product.RestaurantID,
			Title:        product.Title,
			Description:  product.Description,
			CreatedAt:    product.CreatedAt,
			UpdatedAt:    product.UpdatedAt,
		})
	}

	return ctx.JSON(http.StatusOK, products)
}

func registerProduct(ctx echo.Context) error {
	authUser, err := getAuthUser(ctx)
	if err != nil {
		return echo.ErrUnauthorized
	}
	// Only owners can create product
	if authUser.Role != models.RoleOwner {
		return echo.ErrUnauthorized
	}

	payload := struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description" validate:"required"`
	}{}
	if err := ctx.Bind(&payload); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(payload); err != nil {
		return echo.ErrBadRequest
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
	if err := db.Connection.First(restaurant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrUnauthorized
		}
		return echo.ErrInternalServerError
	}

	product := &models.Product{
		RestaurantID: restaurantID,
		Title:        payload.Title,
		Description:  payload.Description,
	}
	tx := db.Connection.Create(product)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrDuplicatedKey) {
			return echo.ErrConflict
		}
		return echo.ErrInternalServerError
	}

	return ctx.JSON(http.StatusCreated, map[string]string{
		"product_id": string(product.ID.String()),
	})
}
