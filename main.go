package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/buivuanh/elotusteam-hackathon/controller"
	"github.com/buivuanh/elotusteam-hackathon/domain"
	"github.com/buivuanh/elotusteam-hackathon/infrastructure/postgresql"
	"github.com/buivuanh/elotusteam-hackathon/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var config *domain.Config

func init() {
	config = &domain.Config{}
	config.Default()
}

func main() {
	// Read configs from env
	config.DBConnectString = os.Getenv("DATABASE_URL")
	config.JWTSigningKey = []byte(os.Getenv("JWT_SIGNING_KEY"))
	config.JWTTokenExpirationHour, _ = strconv.Atoi(os.Getenv("JWT_TOKEN_EXPIRATION_HOUR"))
	config.DataStorePath = os.Getenv("DATA_STORE_PATH")

	// Create database connection pool
	dbPool, err := utils.NewConnectionPool(context.Background(), config.DBConnectString)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer dbPool.Close()

	// Runs migrations
	if err = utils.Up(context.Background(), config.DBConnectString); err != nil {
		fmt.Fprintf(os.Stderr, "migrate up failed: %v\n", err)
		os.Exit(1)
	}

	// Init repo
	userRepo := postgresql.NewUserRepo()
	fileRepo := postgresql.NewFileRepository()

	authen := &Authen{
		DB:       dbPool,
		userRepo: userRepo,
		config:   config,
	}

	// Init http controller
	user := &controller.UserHttpService{
		Authen:   authen,
		DB:       dbPool,
		UserRepo: userRepo,
	}
	fileStore := &controller.FileStoreHttpService{
		DB:            dbPool,
		FileRepo:      fileRepo,
		DataStorePath: config.DataStorePath,
	}

	// Define apis
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(authen.jwtMiddleware)
	e.POST("/login", user.Login)
	e.POST("/register", user.Register)
	e.POST("/upload", fileStore.Upload, authen.checkContentType("image/"), middleware.BodyLimit("8M"))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Start server
	go func() {
		if err = e.Start(config.ServerPort); err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
