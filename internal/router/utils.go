package router

import (
	"errors"

	"github.com/labstack/echo/v4"
)

func getAuthUser(ctx echo.Context) (*authUser, error) {
	authUser, ok := ctx.Get(string(userIDKey)).(authUser)
	if !ok {
		return nil, errors.New("failed to retrieve auth user")
	}
	return &authUser, nil
}
