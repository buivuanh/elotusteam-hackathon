package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/buivuanh/elotusteam-hackathon/domain"
	"github.com/buivuanh/elotusteam-hackathon/infrastructure"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type UserHttpService struct {
	DB     *pgxpool.Pool
	Authen interface {
		GenJWTToken(c echo.Context) (string, error)
	}
	UserRepo infrastructure.UserRepo
}

func (u *UserHttpService) verifyUserAndPassword(ctx context.Context, userName, password string) (*domain.User, error) {
	user, err := u.UserRepo.GetByUserName(ctx, u.DB, userName)
	if err != nil {
		return nil, fmt.Errorf("u.UserRepo.GetByID: %w", err)
	}
	if err = user.VerifyPassword(password); err != nil {
		return nil, fmt.Errorf("VerifyPassword: %w", err)
	}

	return user, nil
}

func (u *UserHttpService) Login(c echo.Context) error {
	userName := c.FormValue("user_name")
	password := c.FormValue("password")

	// Verify user and password
	user, err := u.verifyUserAndPassword(c.Request().Context(), userName, password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("ValidateUserAndPassword: %v", err))
	}

	// Generate jwt token
	c.Set("user_id", user.UserID)
	c.Set("user_name", user.UserName)
	signedToken, err := u.Authen.GenJWTToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("GenJWTToken: %v", err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": signedToken,
	})
}

func (u *UserHttpService) Register(c echo.Context) error {
	userName := c.FormValue("user_name")
	password := c.FormValue("password")

	user := &domain.User{
		UserName: userName,
	}
	if err := user.HashPassword(password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("HashedPassword: %v", err))
	}
	user, err := u.UserRepo.Insert(c.Request().Context(), u.DB, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("UserRepo.Insert: %v", err))
	}

	return c.NoContent(http.StatusCreated)
}
