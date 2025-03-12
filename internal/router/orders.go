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

func bindOrdersRouter(router *echo.Group) {
	group := router.Group("/restaurants/:restaurant_id/products")
	group.POST("", createOrder)
}

type OrderResponse struct {
	OrderID      uuid.UUID          `json:"id"`
	RestaurantID uuid.UUID          `json:"restaurant_id"`
	StaffID      uuid.UUID          `json:"staff_id"`
	TableNumber  string             `json:"table_number"`
	Status       models.OrderStatus `json:"status"`
	TotalAmount  float64            `json:"total_amount"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	Items        struct {
		ItemID    uuid.UUID `json:"item_id"`
		ProductID uuid.UUID `json:"product_id"`
		Quantity  uint32    `json:"quantity"`
	} `json:"items"`
}

func createOrder(ctx echo.Context) error {
	authUser, err := getAuthUser(ctx)
	if err != nil {
		return echo.ErrUnauthorized
	}
	// Only staff can create orders
	if authUser.Role != models.RoleStaff {
		return echo.ErrUnauthorized
	}

	restaurantID, err := uuid.Parse(ctx.Param("restaurant_id"))
	if err != nil {
		return echo.ErrBadRequest
	}

	payload := struct {
		Products []struct {
			ProductID uuid.UUID `json:"product_id" validate:"required"`
			Quantity  uint32    `json:"quantity" validate:"required"`
		} `json:"products" validate:"required"`
		TableNumber string `json:"table_number" validate:"required"`
	}{}
	if err := ctx.Bind(&payload); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(payload); err != nil {
		return echo.ErrBadRequest
	}
	if len(payload.Products) == 0 {
		return echo.ErrBadRequest
	}

	db := ctx.(*routerContext).GetDatabase()

	// check if staff is creating the order for their own restaurant
	if err := db.Connection.
		First(&models.Staff{
			ID:           authUser.UserID,
			RestaurantID: restaurantID,
		}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrUnauthorized
		}
		return echo.ErrInternalServerError
	}

	orderItems := make([]models.OrderItem, 0, len(payload.Products))
	for _, product := range payload.Products {
		orderItems = append(orderItems, models.OrderItem{
			ProductID: product.ProductID,
			Quantity:  product.Quantity,
		})
	}

	order := &models.Order{
		RestaurantID: restaurantID,
		StaffID:      authUser.UserID,
		TableNumber:  payload.TableNumber,
		TotalAmount:  models.CalculateTotalAmount(orderItems),
		OrderItems:   orderItems,
	}

	if err := db.Connection.Create(order).Error; err != nil {
		return echo.ErrInternalServerError
	}

	return ctx.JSON(http.StatusCreated, map[string]string{
		"order_id": order.ID.String(),
	})
}
