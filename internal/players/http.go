package players

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/raphael-trzpit/fifa-test/internal/auth"
	uuid "github.com/satori/go.uuid"
)

// Handler contains the dependencies for all the http handlers.
type Handle struct {
	Repository PlayerRepository
}

// GetAllPlayers will return all players from the current user's team
func (h *Handle) GetAllPlayers(w http.ResponseWriter, r *http.Request) {
	currentUser := auth.CurrentUserFromContext(r.Context())
	if currentUser == nil {
		http.Error(w, "no current user", http.StatusBadRequest)
		return
	}

	players, err := h.Repository.GetAllByTeamID(currentUser.TeamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(players)
}

// GetAllPlayers will return one player by its id.
// If this player is not in the current user's team, it returns an error
// It needs the use of httprouter and an id named param.
func (h *Handle) GetPlayerByID(w http.ResponseWriter, r *http.Request) {
	currentUser := auth.CurrentUserFromContext(r.Context())
	if currentUser == nil {
		http.Error(w, "no current user", http.StatusBadRequest)
		return
	}

	playerUuid, err := uuid.FromString(httprouter.ParamsFromContext(r.Context()).ByName("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	player, err := h.Repository.GetByID(playerUuid)
	if err != nil {
		if errors.As(err, PlayerNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if player.TeamID != currentUser.TeamID {
		http.Error(w, "this player doesn't belong to your team", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(player)
}

// CreatePlayer will create a new player in the current user's team, and generate an id for him.
func (h *Handle) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	currentUser := auth.CurrentUserFromContext(r.Context())
	if currentUser == nil {
		http.Error(w, "no current user", http.StatusBadRequest)
		return
	}

	var payload  struct {
		FirstName string
		LastName string
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	player := &Player{
		ID:        uuid.NewV4(),
		TeamID:    currentUser.TeamID,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
	}
	if err := h.Repository.Create(player); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(player)
}

// UpdatePlayer will update one player by its id.
// If this player is not in the current user's team, it returns an error
// It needs the use of httprouter and an id named param.
func (h *Handle) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	currentUser := auth.CurrentUserFromContext(r.Context())
	if currentUser == nil {
		http.Error(w, "no current user", http.StatusBadRequest)
		return
	}

	playerUuid, err := uuid.FromString(httprouter.ParamsFromContext(r.Context()).ByName("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var payload  struct {
		FirstName string
		LastName string
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	player, err := h.Repository.GetByID(playerUuid)
	if err != nil {
		if errors.As(err, PlayerNotFound) {
			http.NotFound(w, r)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if player.TeamID != currentUser.TeamID {
		http.Error(w, "this player doesn't belong to your team", http.StatusForbidden)
		return
	}

	player.FirstName = payload.FirstName
	player.LastName = payload.LastName

	if err := h.Repository.Update(player); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(player)
}

// DeletePlayer will delete one player by its id.
// If this player is not in the current user's team, it returns an error
// It needs the use of httprouter and an id named param.
func (h *Handle) DeletePlayer(w http.ResponseWriter, r *http.Request) {
	currentUser := auth.CurrentUserFromContext(r.Context())
	if currentUser == nil {
		http.Error(w, "no current user", http.StatusBadRequest)
		return
	}

	playerUuid, err := uuid.FromString(httprouter.ParamsFromContext(r.Context()).ByName("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	player, err := h.Repository.GetByID(playerUuid)
	if err != nil {
		if errors.As(err, PlayerNotFound) {
			http.NotFound(w, r)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if player.TeamID != currentUser.TeamID {
		http.Error(w, "this player doesn't belong to your team", http.StatusForbidden)
		return
	}

	if err := h.Repository.Delete(player.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
