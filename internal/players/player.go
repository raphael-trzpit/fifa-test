package players

import (
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// Player is a football player that belongs to a team.
type Player struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	TeamID    uuid.UUID
	FirstName string
	LastName  string
}

var ErrPlayerNotFound = errors.New("player not found")

// PlayerRepository is a service which will handle the storage of the players.
type PlayerRepository interface {
	Create(*Player) error
	Update(*Player) error
	Delete(uuid.UUID) error
	GetByID(uuid.UUID) (*Player, error)
	GetAllByTeamID(uuid.UUID) ([]*Player, error)
}
