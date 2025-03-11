package router

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/roushou/pocpoc/internal/database"
)

type routerContext struct {
	echo.Context
	database *database.Database
}

func (ctx *routerContext) GetDatabase() *database.Database {
	return ctx.database
}

// withRouterContext extends echo.Context by setting up Services into it.
//
// IMPORTANT: This middleware should be called before any other middlewares and routers.
func withRouterContext(database *database.Database) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			rc := &routerContext{Context: ctx, database: database}
			return next(rc)
		}
	}
}

func NewRouter(database *database.Database) *echo.Echo {
	router := echo.New()
	router.Validator = &Validator{validator: validator.New()}

	group := router.Group("/api")

	// Middlewares
	group.Use(withRouterContext(database)) // !!! This middleware should be called before anything else
	group.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// TODO: Load from configuration
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
	}))
	group.Use(middleware.Logger())
	group.Use(middleware.Recover())

	// Routers
	bindHealthRouter(group)
	bindAuthRouter(group)

	// Need auth
	restricted := group.Group("")
	restricted.Use(AuthMiddleware(jwtSecretKey))
	bindRestaurantsRouter(restricted)

	return router
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
