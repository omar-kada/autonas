// Package middlewares provides HTTP middleware functionality.
package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"omar-kada/autonas/api"
	"omar-kada/autonas/internal/users"
	"omar-kada/autonas/models"
)

type contextKey string

const (
	_tokenKey                   = "token"
	_refreshTokenKey            = "refreshToken"
	_usernameKey     contextKey = "username"
)

// ContextWithUsername adds user information to the context.
// @param ctx context.Context - the context to add user information to
// @param user models.User - the user information to add
// @return context.Context - the context with user information added
func ContextWithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, _usernameKey, username)
}

// UsernameFromContext retrieves user information from the context.
// @param ctx context.Context - the context to retrieve user information from
// @return models.User - the user information retrieved
// @return bool - true if user information was found, false otherwise
func UsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(_usernameKey).(string)
	return username, ok
}

// AuthMiddleware provides authentication middleware.
// @param next http.Handler - the next handler in the chain
// @param authService user.AuthService - the authentication service
// @return http.Handler - the authentication middleware
func AuthMiddleware(next http.Handler, authService users.AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url, ok := strings.CutPrefix(r.URL.Path, "/api/")
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		switch url {
		case "auth/register":
			registerHandler(w, r, authService)
			return
		case "auth/login":
			loginHandler(w, r, authService)
			return
		case "auth/logout":
			logoutHandler(w, r, authService)
			return
		case "auth/refresh":
			refreshHandler(w, r, authService)
			return
		}
		inWhiteList := isWhitelisted(url, r.Method)

		username, err := getUsernameFromCookies(r, authService)
		if err != nil {
			slog.Error(err.Error())
			if !inWhiteList {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
		}
		r = r.WithContext(ContextWithUsername(r.Context(), username))

		next.ServeHTTP(w, r)
	})
}

func getUsernameFromCookies(r *http.Request, authService users.AuthService) (string, error) {
	token := getTokenFromCookies(r)
	if token.RefreshToken == "" {
		return "", errors.New("no auth available")
	}
	// username, err := authService.GetUsernameByToken(token)
	// if err == nil {
	// 	return username, nil
	// }
	// token, err = authService.RefreshToken(token)
	// if err != nil {
	// 	return "", err
	// }
	// setTokenInCookies(w, token)
	return authService.GetUsernameByToken(token)
}

var _whitelisted = map[string][]string{
	"user": {"GET"},
}

func isWhitelisted(url, method string) bool {
	if methods, ok := _whitelisted[url]; ok {
		return slices.Contains(methods, method)
	}
	return false
}

func registerHandler(w http.ResponseWriter, r *http.Request, authService users.AuthService) {
	switch r.Method {
	case http.MethodPost:
		var req api.Credentials

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Error(err.Error())
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Username == "" || req.Password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}

		token, err := authService.Register(models.Credentials{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "Registration failed", http.StatusInternalServerError)
			return
		}

		setTokenInCookies(w, token)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(api.BooleanResponse{
			Success: true,
		})
		return
	case http.MethodGet:
		hasUsers, err := authService.IsRegistered()
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "testing", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(api.AuthAPIRegistered200JSONResponse{
			Registered: hasUsers,
		})
		return
	default:
		http.Error(w, "invalid method", http.StatusUnauthorized)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request, authService users.AuthService) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}

	var req api.Credentials

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	auth, err := authService.Login(models.Credentials{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}

	setTokenInCookies(w, auth)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.BooleanResponse{
		Success: true,
	})
}

func refreshHandler(w http.ResponseWriter, r *http.Request, authService users.AuthService) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}

	token := getTokenFromCookies(r)

	if token.RefreshToken == "" {
		slog.Error("invalid token value")
		http.Error(w, "Logout failed", http.StatusUnauthorized)
		return
	}

	newToken, err := authService.RefreshToken(token)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Refresh failed", http.StatusUnauthorized)
		return
	}

	setTokenInCookies(w, newToken)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.BooleanResponse{
		Success: true,
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request, authService users.AuthService) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}
	_, err := getUsernameFromCookies(r, authService)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Couldn't find user", http.StatusUnauthorized)
		return
	}

	token := getTokenFromCookies(r)

	if token.Value == "" || token.RefreshToken == "" {
		slog.Error("invalid token value")
		http.Error(w, "Logout failed", http.StatusUnauthorized)
		return
	}

	err = authService.Logout(token)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Logout failed", http.StatusUnauthorized)
		return
	}

	setTokenInCookies(w, models.Token{
		Value:          "",
		Expires:        time.Unix(0, 0),
		RefreshToken:   "",
		RefreshExpires: time.Unix(0, 0),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.BooleanResponse{
		Success: true,
	})
}

func getTokenFromCookies(r *http.Request) models.Token {
	cookie, err := r.Cookie(_tokenKey)
	if err != nil {
		cookie = &http.Cookie{
			Value: "",
		}
	}
	refreshCookie, err := r.Cookie(_refreshTokenKey)
	if err != nil {
		refreshCookie = &http.Cookie{
			Value: "",
		}
	}
	return models.Token{
		Value:          models.TokenValue(cookie.Value),
		Expires:        cookie.Expires,
		RefreshToken:   models.TokenValue(refreshCookie.Value),
		RefreshExpires: refreshCookie.Expires,
	}
}

func setTokenInCookies(w http.ResponseWriter, token models.Token) {
	http.SetCookie(w, &http.Cookie{
		Name:     _tokenKey,
		Value:    string(token.Value),
		Expires:  token.Expires,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/api",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     _refreshTokenKey,
		Value:    string(token.RefreshToken),
		Expires:  token.RefreshExpires,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/api",
	})
}
