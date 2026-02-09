package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"omar-kada/autonas/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of user.AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(credentials models.Credentials) (models.Auth, error) {
	args := m.Called(credentials)
	return args.Get(0).(models.Auth), args.Error(1)
}

func (m *MockAuthService) Register(credentials models.Credentials) (models.Auth, error) {
	args := m.Called(credentials)
	return args.Get(0).(models.Auth), args.Error(1)
}

func (m *MockAuthService) IsRegistered() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthService) Logout(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockAuthService) GetUserByToken(token string) (models.User, error) {
	args := m.Called(token)
	return args.Get(0).(models.User), args.Error(1)
}

func TestAuthMiddleware_Register(t *testing.T) {
	mockAuthService := new(MockAuthService)
	expectedAuth := models.Auth{
		Token:     "testtoken",
		ExpiresIn: time.Now().Add(24 * time.Hour),
	}

	mockAuthService.On("Register", models.Credentials{
		Username: "testuser",
		Password: "testpass",
	}).Return(expectedAuth, nil)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), mockAuthService)

	reqBody := `{"username":"testuser","password":"testpass"}`
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_RegisterGet(t *testing.T) {
	mockAuthService := new(MockAuthService)

	// Mock the IsRegistered method to return true
	mockAuthService.On("IsRegistered").Return(true, nil)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), mockAuthService)

	req := httptest.NewRequest("GET", "/api/register", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"registered": true}`, rr.Body.String())
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_Login(t *testing.T) {
	mockAuthService := new(MockAuthService)
	expectedAuth := models.Auth{
		Token:     "testtoken",
		ExpiresIn: time.Now().Add(24 * time.Hour),
	}

	mockAuthService.On("Login", models.Credentials{
		Username: "testuser",
		Password: "testpass",
	}).Return(expectedAuth, nil)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), mockAuthService)

	reqBody := `{"username":"testuser","password":"testpass"}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_Logout(t *testing.T) {
	mockAuthService := new(MockAuthService)

	mockAuthService.On("Logout", "testtoken").Return(nil)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), mockAuthService)

	req := httptest.NewRequest("POST", "/api/logout", http.NoBody)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "testtoken",
	})
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_AuthorizedAccess(t *testing.T) {
	mockAuthService := new(MockAuthService)
	expectedUser := models.User{
		Username: "testuser",
	}

	mockAuthService.On("GetUserByToken", "testtoken").Return(expectedUser, nil)

	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, "testuser", user.Username)
		w.WriteHeader(http.StatusOK)
	}), mockAuthService)

	req := httptest.NewRequest("GET", "/api/protected", http.NoBody)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "testtoken",
	})
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_UnauthorizedAccess(t *testing.T) {
	mockAuthService := new(MockAuthService)

	mockAuthService.On("GetUserByToken", "invalidtoken").Return(models.User{}, errors.New("invalid token"))

	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), mockAuthService)

	req := httptest.NewRequest("GET", "/api/protected", http.NoBody)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "invalidtoken",
	})
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_UnauthorizedEmptyToken(t *testing.T) {
	mockAuthService := new(MockAuthService)

	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), mockAuthService)

	req := httptest.NewRequest("GET", "/api/protected", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_WhitelistedAccess(t *testing.T) {
	mockAuthService := new(MockAuthService)

	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), mockAuthService)

	req := httptest.NewRequest("GET", "/api/user", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_RegisterInvalidRequestBody(t *testing.T) {
	mockAuthService := new(MockAuthService)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), mockAuthService)

	reqBody := `{"username":"testuser"}` // missing password
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_RegisterMissingCredentials(t *testing.T) {
	mockAuthService := new(MockAuthService)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), mockAuthService)

	reqBody := `{"username":"","password":""}` // empty username and password
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_RegisterFailure(t *testing.T) {
	mockAuthService := new(MockAuthService)

	mockAuthService.On("Register", models.Credentials{
		Username: "testuser",
		Password: "testpass",
	}).Return(models.Auth{}, errors.New("registration failed"))

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), mockAuthService)

	reqBody := `{"username":"testuser","password":"testpass"}`
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_LoginInvalidMethod(t *testing.T) {
	mockAuthService := new(MockAuthService)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), mockAuthService)

	req := httptest.NewRequest("GET", "/api/login", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_LoginInvalidRequestBody(t *testing.T) {
	mockAuthService := new(MockAuthService)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), mockAuthService)

	reqBody := `{"username":"testuser"}` // missing password
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_LoginMissingCredentials(t *testing.T) {
	mockAuthService := new(MockAuthService)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), mockAuthService)

	reqBody := `{"username":"","password":""}` // empty username and password
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_LoginFailure(t *testing.T) {
	mockAuthService := new(MockAuthService)

	mockAuthService.On("Login", models.Credentials{
		Username: "testuser",
		Password: "testpass",
	}).Return(models.Auth{}, errors.New("login failed"))

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), mockAuthService)

	reqBody := `{"username":"testuser","password":"testpass"}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_LogoutInvalidMethod(t *testing.T) {
	mockAuthService := new(MockAuthService)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), mockAuthService)

	req := httptest.NewRequest("GET", "/api/logout", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_LogoutMissingToken(t *testing.T) {
	mockAuthService := new(MockAuthService)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), mockAuthService)

	req := httptest.NewRequest("POST", "/api/logout", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	mockAuthService.AssertExpectations(t)
}

func TestAuthMiddleware_LogoutFailure(t *testing.T) {
	mockAuthService := new(MockAuthService)

	mockAuthService.On("Logout", "testtoken").Return(errors.New("logout failed"))

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), mockAuthService)

	req := httptest.NewRequest("POST", "/api/logout", http.NoBody)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "testtoken",
	})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	mockAuthService.AssertExpectations(t)
}
