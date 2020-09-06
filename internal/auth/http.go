package auth

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// CreateUserHandler returns an http handler that create a user.
// It will also generate a team id for the created user.
func CreateUserHandler(repo UserRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Name     string
			Password string
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hash, err := HashPassword(payload.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := &User{
			Name:         payload.Name,
			PasswordHash: hash,
			TeamID:       uuid.NewV4(),
		}
		if err := repo.Create(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	})
}

// AuthMiddleware returns an authentication middleware.
// It will check if the user is logged in with basic auth and that the user is present in the user repository.
func AuthMiddleware(repo UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, hasAuth := r.BasicAuth()
			if !hasAuth {
				http.Error(w, "not logged in", http.StatusUnauthorized)
				return
			}

			user, err := repo.GetUserByName(username)
			if err != nil {
				if errors.As(err, UserNotFound.Error()) {
					http.Error(w, "not logged in", http.StatusUnauthorized)
					return
				}

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if ok := CheckPasswordHash(password, user.PasswordHash); !ok {
				http.Error(w, "not logged in", http.StatusUnauthorized)
				return
			}

			r = r.WithContext(ContextWithCurrentUser(r.Context(), user))

			next.ServeHTTP(w, r)
		})
	}
}
