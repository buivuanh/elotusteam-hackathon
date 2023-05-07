package controller

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/buivuanh/elotusteam-hackathon/domain"
	"github.com/buivuanh/elotusteam-hackathon/infrastructure"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	// Create a new instance of the Echo framework
	e := echo.New()

	// Create mock objects
	mockUserRepo := new(infrastructure.MockUserRepo)

	// Create a test user
	user := &domain.User{
		UserName: "testuser",
	}

	// Mock the Insert method of the UserRepo
	mockUserRepo.On("Insert", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		actual := args[2].(*domain.User)
		assert.Equal(t, user.UserName, actual.UserName)
		assert.NotEmpty(t, actual.HashedPassword)
	}).Return(user, nil)

	// Create the UserHttpService instance with the mocks
	userService := &UserHttpService{
		UserRepo: mockUserRepo,
	}

	// Create a new HTTP request context
	req := httptest.NewRequest(http.MethodPost, "/register", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set the form values
	c.Request().Form = url.Values{}
	c.Request().Form.Set("user_name", "testuser")
	c.Request().Form.Set("password", "testpassword")

	// Call the Register method
	err := userService.Register(c)
	assert.NoError(t, err)

	// Assert the response
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Empty(t, rec.Body.String())
}
