package router

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/roushou/pocpoc/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthUser(t *testing.T) {
	t.Run("auth user exists in context", func(t *testing.T) {
		e := echo.New()
		ctx := e.NewContext(nil, nil)

		userID, err := uuid.NewV7()
		assert.NoError(t, err)

		expectedUser := authUser{UserID: userID, Role: models.RoleOwner}
		ctx.Set(string(userIDKey), expectedUser)

		authUser, err := getAuthUser(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser.UserID, authUser.UserID)
		assert.Equal(t, expectedUser.Role, authUser.Role)
	})

	t.Run("invalid auth user type", func(t *testing.T) {
		e := echo.New()
		ctx := e.NewContext(nil, nil)
		ctx.Set(string(userIDKey), "not-an-auth-user")

		authUser, err := getAuthUser(ctx)

		assert.Error(t, err)
		assert.Nil(t, authUser)
		assert.Equal(t, errors.New("failed to retrieve auth user"), err)
	})

	t.Run("no auth user in context", func(t *testing.T) {
		e := echo.New()
		ctx := e.NewContext(nil, nil)
		ctx.Set(string(userIDKey), nil)

		authUser, err := getAuthUser(ctx)

		assert.Error(t, err)
		assert.Nil(t, authUser)
		assert.Equal(t, errors.New("failed to retrieve auth user"), err)
	})
}

