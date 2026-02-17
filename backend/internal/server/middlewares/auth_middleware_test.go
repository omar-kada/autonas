package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/internal/users"
	"omar-kada/autonas/models"
	"omar-kada/autonas/testutil"

	"github.com/stretchr/testify/assert"
)

var userCreds = models.Credentials{
	Username: "username",
	Password: "password",
}

func newUsersService(t *testing.T) users.Service {
	store, err := storage.NewUsersStorage(testutil.NewMemoryStorage())
	assert.NoError(t, err)
	return users.NewService(store)
}

func withInitUsers(t *testing.T, userService users.Service, creds models.Credentials) (users.Service, models.Token) {
	token, err := userService.Register(creds)
	assert.NoError(t, err)

	return userService, token
}

func TestAuthMiddleware_Register_RealService(t *testing.T) {
	userService := newUsersService(t)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), userService)

	reqBody := `{"username":"testuser","password":"testpass"}`
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	checkCookiesAreNot(t, rr, "", "")
}

func TestAuthMiddleware_RegisterGet_RealService(t *testing.T) {
	userService, _ := withInitUsers(t, newUsersService(t), userCreds)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), userService)

	req := httptest.NewRequest("GET", "/api/register", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"registered": true}`, rr.Body.String())
	checkCookiesAre(t, rr, "", "")
}

func TestAuthMiddleware_Login_RealService(t *testing.T) {
	userService := newUsersService(t)
	userService, token := withInitUsers(t, userService, userCreds)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), userService)

	reqBody := `{"username":"username","password":"password"}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	checkCookiesAreNot(t, rr, string(token.Value), string(token.RefreshToken))
	checkCookiesAreNot(t, rr, "", "")
}

func TestAuthMiddleware_Logout_RealService(t *testing.T) {
	userService := newUsersService(t)
	userService, token := withInitUsers(t, userService, userCreds)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), userService)

	req := httptest.NewRequest("POST", "/api/logout", http.NoBody)
	req.AddCookie(&http.Cookie{
		Name:    _tokenKey,
		Value:   string(token.Value),
		Expires: token.Expires,
	})
	req.AddCookie(&http.Cookie{
		Name:    _refreshTokenKey,
		Value:   string(token.RefreshToken),
		Expires: token.RefreshExpires,
	})
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	checkCookiesAre(t, rr, "", "")
}

func TestAuthMiddleware_AuthorizedAccess_RealService(t *testing.T) {
	userService := newUsersService(t)
	userService, token := withInitUsers(t, userService, userCreds)

	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, ok := UsernameFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, "username", username)
		w.WriteHeader(http.StatusOK)
	}), userService)

	req := httptest.NewRequest("GET", "/api/protected", http.NoBody)
	req.AddCookie(&http.Cookie{
		Name:  _tokenKey,
		Value: string(token.Value),
	})
	req.AddCookie(&http.Cookie{
		Name:  _refreshTokenKey,
		Value: string(token.RefreshToken),
	})
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_WhitelistedAccess_RealService(t *testing.T) {
	userService := newUsersService(t)

	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), userService)

	req := httptest.NewRequest("GET", "/api/user", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	checkCookiesAre(t, rr, "", "")
}

func TestAuthMiddleware_RegisterInvalidRequestBody_RealService(t *testing.T) {
	userService := newUsersService(t)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), userService)

	reqBody := `{"username":"testuser"}` // missing password
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	checkCookiesAre(t, rr, "", "")
}

func TestAuthMiddleware_RegisterMissingCredentials_RealService(t *testing.T) {
	userService := newUsersService(t)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), userService)

	reqBody := `{"username":"","password":""}` // empty username and password
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	checkCookiesAre(t, rr, "", "")
}

func TestAuthMiddleware_RegisterFailure_RealService(t *testing.T) {
	userService := newUsersService(t)
	// Register first user to prevent registration
	_, _ = withInitUsers(t, userService, userCreds)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), userService)

	reqBody := `{"username":"testuser","password":"testpass"}`
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	checkCookiesAre(t, rr, "", "")
}

func TestAuthMiddleware_LoginInvalidMethod_RealService(t *testing.T) {
	userService := newUsersService(t)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), userService)

	req := httptest.NewRequest("GET", "/api/login", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	checkCookiesAre(t, rr, "", "")
}

func TestAuthMiddleware_LoginInvalidRequestBody_RealService(t *testing.T) {
	userService := newUsersService(t)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), userService)

	reqBody := `{"username":"testuser"}` // missing password
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	checkCookiesAre(t, rr, "", "")
}

func TestAuthMiddleware_LoginMissingCredentials_RealService(t *testing.T) {
	userService := newUsersService(t)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), userService)

	reqBody := `{"username":"","password":""}` // empty username and password
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	checkCookiesAre(t, rr, "", "")
}

func TestAuthMiddleware_LoginFailure_RealService(t *testing.T) {
	userService := newUsersService(t)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), userService)

	reqBody := `{"username":"wronguser","password":"wrongpass"}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	checkCookiesAre(t, rr, "", "")
}

func TestAuthMiddleware_LogoutInvalidMethod_RealService(t *testing.T) {
	userService := newUsersService(t)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), userService)

	req := httptest.NewRequest("GET", "/api/logout", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	checkCookiesAre(t, rr, "", "")
}

func TestAuthMiddleware_LogoutMissingToken_RealService(t *testing.T) {
	userService := newUsersService(t)

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), userService)

	req := httptest.NewRequest("POST", "/api/logout", http.NoBody)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	checkCookiesAre(t, rr, "", "")
}

func TestAuthMiddleware_LogoutFailure_RealService(t *testing.T) {
	userService := newUsersService(t)
	userService, token := withInitUsers(t, userService, userCreds)

	// Invalidate the token by modifying it
	invalidToken := token
	invalidToken.Value = "invalidtoken"

	handler := AuthMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fail() // shouldn't be called
	}), userService)

	req := httptest.NewRequest("POST", "/api/logout", http.NoBody)
	req.AddCookie(&http.Cookie{
		Name:  _tokenKey,
		Value: string(invalidToken.Value),
	})
	req.AddCookie(&http.Cookie{
		Name:  _refreshTokenKey,
		Value: string(invalidToken.RefreshToken),
	})
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func checkCookiesAreNot(t *testing.T, rr *httptest.ResponseRecorder, expectedToken, expectedRefreshToken string) {
	cookies := rr.Result().Cookies()
	var tokenCookie, refreshTokenCookie http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == _tokenKey {
			tokenCookie = *cookie
		}
		if cookie.Name == _refreshTokenKey {
			refreshTokenCookie = *cookie
		}
	}

	assert.NotNil(t, tokenCookie, "Token cookie should be set")
	assert.NotEqual(t, expectedToken, tokenCookie.Value, "Token cookie value should match")

	assert.NotNil(t, refreshTokenCookie, "Refresh token cookie should be set")
	assert.NotEqual(t, expectedRefreshToken, refreshTokenCookie.Value, "Refresh token cookie value should match")
}

func checkCookiesAre(t *testing.T, rr *httptest.ResponseRecorder, expectedToken, expectedRefreshToken string) {
	cookies := rr.Result().Cookies()
	var tokenCookie, refreshTokenCookie http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == _tokenKey {
			tokenCookie = *cookie
		}
		if cookie.Name == _refreshTokenKey {
			refreshTokenCookie = *cookie
		}
	}

	assert.NotNil(t, tokenCookie, "Token cookie should be set")
	assert.Equal(t, expectedToken, tokenCookie.Value, "Token cookie value should match")

	assert.NotNil(t, refreshTokenCookie, "Refresh token cookie should be set")
	assert.Equal(t, expectedRefreshToken, refreshTokenCookie.Value, "Refresh token cookie value should match")
}
