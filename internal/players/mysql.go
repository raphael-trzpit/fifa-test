package players

import (
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// PlayerRepositoryMysl is an implementation of PlayerRepository using a Mysl database as storage.
type PlayerRepositoryMysl struct {
	db *gorm.DB
}

func NewPlayerRepositoryMysl(db *gorm.DB) (*PlayerRepositoryMysl, error) {
	if db == nil {
		return nil, errors.New("cannot instantiate new player repository Mysl: no db provided")
	}

	return &PlayerRepositoryMysl{db: db}, nil
}

func (r *PlayerRepositoryMysl) Create(player *Player) error {
	return errors.Wrap(r.db.Create(player).Error, "cannot save player")
}

func (r *PlayerRepositoryMysl) Update(player *Player) error {
	return errors.Wrap(r.db.Save(player).Error, "cannot save player")
}

func (r *PlayerRepositoryMysl) Delete(uuid uuid.UUID) error {
	return errors.Wrap(r.db.Delete(&Player{}, uuid).Error, "cannot delete player")
}

func (r *PlayerRepositoryMysl) GetByID(uuid uuid.UUID) (*Player, error) {
	player :=  &Player{}
	err := r.db.First(player, uuid).Error

	if err != nil {
		if errors.As(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(err, PlayerNotFound.Error())
		}
		return nil, errors.Wrap(err, "cannot get player")
	}

	return player, nil
}

func (r *PlayerRepositoryMysl) GetAllByTeamID(uuid uuid.UUID) ([]*Player, error) {
	var players []*Player
	err := r.db.Where("team_id = ?", uuid).Find(&players).Error
	if err != nil {
		if errors.As(err, gorm.ErrRecordNotFound) {
			return players, nil
		}
		return nil, errors.Wrap(err, "cannot get players")
	}

	return players, nil
}
