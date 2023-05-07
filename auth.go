package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/buivuanh/elotusteam-hackathon/domain"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
)

var ignoreAuthEndpoint = []string{
	"/register",
	"/login",
	"/",
}

type Authen struct {
	DB       *pgxpool.Pool
	userRepo interface {
		GetByID(ctx context.Context, db *pgxpool.Pool, userID int) (*domain.User, error)
	}
	config *domain.Config
}

type Claims struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

func (a *Authen) skipValidate(path string) bool {
	if slices.Contains(ignoreAuthEndpoint, path) {
		return true
	}

	return false
}

func (a *Authen) checkContentType(expectContentType string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			file, err := c.FormFile("file")
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Failed to retrieve uploaded file")
			}

			src, err := file.Open()
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Failed to open uploaded file")
			}
			defer src.Close()

			// Read the first 512 bytes of the file to determine its content type
			buf := make([]byte, 512)
			_, err = io.ReadFull(src, buf)
			if err != nil && err != io.EOF {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read uploaded file")
			}

			// check content type
			contentType := http.DetectContentType(buf)
			if !strings.HasPrefix(contentType, expectContentType) {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Uploaded file is not an %s", expectContentType))
			}
			c.Set("ContentType", contentType)

			return next(c)
		}
	}
}

// jwtMiddleware will verify jwt token of request and
// validate user info in db
func (a *Authen) jwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if a.skipValidate(c.Path()) {
			return next(c)
		}
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.ErrUnauthorized
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return a.config.JWTSigningKey, nil
		})

		if err != nil || !token.Valid {
			return echo.ErrUnauthorized
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			return echo.ErrUnauthorized
		}

		c.Set("user_name", claims.UserName)
		c.Set("user_id", claims.UserID)

		if err = a.checkUserInDB(c.Request().Context(), claims.UserID, claims.UserName); err != nil {
			log.Println(err)
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

func (a *Authen) checkUserInDB(ctx context.Context, userID int, userName string) error {
	u, err := a.userRepo.GetByID(ctx, a.DB, userID)
	if err != nil {
		return fmt.Errorf("could not found user %s: %w", userName, err)
	}
	if u.UserName != userName {
		return fmt.Errorf("miss match user name: got %s but existing %s", u.UserName, userName)
	}

	if u.DeletedAt != nil {
		return fmt.Errorf("user was deleted %s", userName)
	}

	return nil
}

// GenJWTToken will return jwt token with HS256 method
func (a *Authen) GenJWTToken(c echo.Context) (string, error) {
	userName := c.Get("user_name").(string)
	userID := c.Get("user_id").(int)
	claims := &Claims{
		UserName: userName,
		UserID:   userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(a.config.JWTTokenExpirationHour)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(a.config.JWTSigningKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
