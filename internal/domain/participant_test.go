package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParticipant(t *testing.T) {
	userID := "user123"
	name := "Test User"
	pictureURL := "http://example.com/pic.jpg"
	identity := Seer
	assert := assert.New(t)

	p := NewParticipant(userID, name, pictureURL, identity)

	assert.Equal(userID, p.UserID, "UID mismatch")
	assert.Equal(name, p.Name, "Name mismatch")
	assert.Equal(pictureURL, p.PictureURL, "PictureURL mismatch")
	assert.Equal(identity, p.Identity, "Identity mismatch")
}

func TestIdentity_String(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		identity Identity
		want     string
	}{
		{WerewolfKing, "狼王"},
		{WhiteWerewolf, "白狼王"},
		{GhostRider, "惡靈騎士"},
		{Werewolf, "狼人"},
		{Seer, "預言家"},
		{Witch, "女巫"},
		{Hunter, "獵人"},
		{Guard, "守衛"},
		{Knight, "騎士"},
		{Magician, "魔術師"},
		{Villager, "平民"},
		{Identity(99), "unknown"},   // Test unknown identity
		{WerewolfBeauty, "unknown"}, // As per current String() implementation
	}

	for _, tt := range tests {
		assert.Equal(tt.want, tt.identity.String(), "Identity(%d).String() mismatch", tt.identity)
	}
}
