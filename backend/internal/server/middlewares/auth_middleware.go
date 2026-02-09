// Package middlewares provides HTTP middleware functionality.
package middlewares

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"omar-kada/autonas/api"
	"omar-kada/autonas/internal/user"
	"omar-kada/autonas/models"
)

type contextKey string

const (
	_tokenKey            = "token"
	_userKey  contextKey = "user"
)

// ContextWithUser adds user information to the context.
// @param ctx context.Context - the context to add user information to
// @param user models.User - the user information to add
// @return context.Context - the context with user information added
func ContextWithUser(ctx context.Context, user models.User) context.Context {
	return context.WithValue(ctx, _userKey, user)
}

// UserFromContext retrieves user information from the context.
// @param ctx context.Context - the context to retrieve user information from
// @return models.User - the user information retrieved
// @return bool - true if user information was found, false otherwise
func UserFromContext(ctx context.Context) (models.User, bool) {
	user, ok := ctx.Value(_userKey).(models.User)
	return user, ok
}

// AuthMiddleware provides authentication middleware.
// @param next http.Handler - the next handler in the chain
// @param authService user.AuthService - the authentication service
// @return http.Handler - the authentication middleware
func AuthMiddleware(next http.Handler, authService user.AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url, ok := strings.CutPrefix(r.URL.Path, "/api/")
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		switch url {
		case "register":
			registerHandler(w, r, authService)
			return
		case "login":
			loginHandler(w, r, authService)
			return
		case "logout":
			logoutHandler(w, r, authService)
			return
		}

		cookie, err := r.Cookie(_tokenKey)
		if err != nil || cookie.Value == "" {
			if !isWhitelisted(url, r.Method) {
				slog.Error(err.Error())
				http.Error(w, "No auth info found", http.StatusUnauthorized)
				return
			}
		} else {

			user, err := authService.GetUserByToken(cookie.Value)
			if err != nil {
				slog.Error(err.Error())
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
			r = r.WithContext(ContextWithUser(r.Context(), user))
		}

		next.ServeHTTP(w, r)
	})
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

func registerHandler(w http.ResponseWriter, r *http.Request, authService user.AuthService) {
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

		auth, err := authService.Register(models.Credentials{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "Registration failed", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     _tokenKey,
			Value:    auth.Token,
			Expires:  auth.ExpiresIn,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Secure:   true,
		})

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
		json.NewEncoder(w).Encode(api.RegisterAPIRegistered200JSONResponse{
			Registered: hasUsers,
		})
		return
	default:
		http.Error(w, "invalid method", http.StatusUnauthorized)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request, authService user.AuthService) {
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

	http.SetCookie(w, &http.Cookie{
		Name:     _tokenKey,
		Value:    auth.Token,
		Expires:  auth.ExpiresIn,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.BooleanResponse{
		Success: true,
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request, authService user.AuthService) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie(_tokenKey)
	if err != nil || cookie.Value == "" {
		slog.Error(err.Error())
		http.Error(w, "Logout failed", http.StatusUnauthorized)
		return
	}

	err = authService.Logout(cookie.Value)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Logout failed", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     _tokenKey,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.BooleanResponse{
		Success: true,
	})
}
