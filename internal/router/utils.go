package router

import (
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func getUserID(ctx echo.Context) (uuid.UUID, error) {
	userID, ok := ctx.Get(string(userIDKey)).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid user UUID")
	}
	return userID, nil
}
