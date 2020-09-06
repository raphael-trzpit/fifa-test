package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckPasswordHash(t *testing.T) {
	tests := []struct {
		name string
		initialPassword, checkedPassword string
		match                            bool
	}{
		{"matching password", "password", "password", true},
		{"not matching passwod", "password", "aze", false},
		{"not matching passwod with uppercase", "password", "Password", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hash, err := HashPassword(test.initialPassword)
			assert.Nil(t, err)

			ok := CheckPasswordHash(test.checkedPassword, hash)
			assert.Equal(t, test.match, ok)
		})
	}
}
