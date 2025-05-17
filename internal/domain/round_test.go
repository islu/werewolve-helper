package domain

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRound(t *testing.T) {
	userID := "owner123"
	inviteNo := "inviteABC"
	round := NewRound(userID, inviteNo)
	assert := assert.New(t) // Use assert.New(t) for convenience

	assert.Equal(userID, round.OwnerID, "OwnerID should match")
	assert.Equal(inviteNo, round.InviteNo, "InviteNo should match")
	assert.Empty(round.Identities, "Identities should be empty initially")
	assert.Empty(round.Participants, "Participants should be empty initially")
	assert.False(round.CreatedAt.IsZero(), "CreateTime should be set")
	assert.False(round.ExpiredAt.IsZero(), "ExpireTime should be set")
	assert.InDelta(2.0, round.ExpiredAt.Sub(round.CreatedAt).Hours(), 0.1, "ExpireTime should be around 2 hours after CreateTime")
	assert.False(round.TempIdentityFlag, "TempIdentityFlag should be false initially")
}

func TestRound_SetIdentity(t *testing.T) {
	ownerID := "owner123"
	nonOwnerID := "nonOwner456"
	round := NewRound(ownerID, "testInvite")
	assert := assert.New(t) // Use assert.New(t) for convenience

	// Test setting identity by owner
	round.SetIdentity(ownerID, Werewolf, 2)
	assert.Len(round.Identities, 2, "Expected 2 identities")
	werewolfCount := 0
	for _, identity := range round.Identities {
		if identity == Werewolf {
			werewolfCount++
		}
	}
	assert.Equal(2, werewolfCount, "Expected 2 Werewolf identities")

	// Test setting identity by non-owner
	round.SetIdentity(nonOwnerID, Villager, 1)
	assert.Len(round.Identities, 2, "Identities should not change after non-owner attempt")
	villagerCount := 0
	for _, identity := range round.Identities {
		if identity == Villager {
			villagerCount++
		}
	}
	assert.Equal(0, villagerCount, "Expected 0 Villager identities after non-owner attempt")

	// Test shuffling (hard to test deterministically without more control over rng)
	// We can at least check if the function runs without error and identities are still present.
	round.SetIdentity(ownerID, Seer, 1)
	assert.Len(round.Identities, 3, "Expected 3 identities after adding Seer")
}

func TestRound_Register(t *testing.T) {
	round := NewRound("owner123", "testInvite")
	round.SetIdentity("owner123", Werewolf, 1)
	round.SetIdentity("owner123", Villager, 1)
	assert := assert.New(t)

	// Test successful registration
	identityStr1 := round.Register("user1", "User One", "url1")
	assert.NotEmpty(identityStr1, "Expected successful registration for user1, got empty string")
	assert.Len(round.Participants, 1, "Expected 1 participant")
	if len(round.Participants) == 1 { // Guard against panic if previous assert fails
		assert.Equal("user1", round.Participants[0].UserID, "Expected participant UID to be user1")
		assert.Equal(identityStr1, round.Participants[0].Identity.String(), "Participant identity string should match returned string")
	}

	identityStr2 := round.Register("user2", "User Two", "url2")
	assert.NotEmpty(identityStr2, "Expected successful registration for user2, got empty string")
	assert.Len(round.Participants, 2, "Expected 2 participants")

	// Test registration when closed
	assert.True(round.IsRegistrationClose(), "Expected registration to be closed")
	identityStr3 := round.Register("user3", "User Three", "url3")
	assert.Empty(identityStr3, "Expected registration to be closed and return empty string")
	assert.Len(round.Participants, 2, "Participants count should not change after closed registration attempt")
}

func TestRound_Again(t *testing.T) {
	round := NewRound("owner123", "testInvite")
	round.SetIdentity("owner123", Werewolf, 1)
	round.SetIdentity("owner123", Villager, 1)
	round.Register("user1", "User One", "url1")
	round.Register("user2", "User Two", "url2")
	assert := assert.New(t)

	originalExpireTime := round.ExpiredAt
	originalIdentities := make([]Identity, len(round.Identities))
	copy(originalIdentities, round.Identities)

	round.Again()

	assert.Empty(round.Participants, "Participants should be empty after Again()")
	assert.True(round.ExpiredAt.After(originalExpireTime), "ExpireTime should be extended")
	assert.ElementsMatch(originalIdentities, round.Identities, "Identities should remain the same (though possibly reordered)")
}

func TestRound_GetParticipantsInfoReplyMessage(t *testing.T) {
	ownerID := "owner123"
	nonOwnerID := "nonOwner456"
	round := NewRound(ownerID, "testInvite")
	round.SetIdentity(ownerID, Werewolf, 1)
	round.SetIdentity(ownerID, Villager, 1)
	round.Register("user1", "User One", "url1")
	assert := assert.New(t)

	// Test by owner
	info := round.GetParticipantsInfoReplyMessage(ownerID)
	expectedPrefix := "目前參與人數: 1/2"
	assert.True(strings.HasPrefix(info, expectedPrefix), "Info should start with '%s', got '%s'", expectedPrefix, info)
	if len(round.Participants) > 0 {
		assert.Contains(info, "User One:"+round.Participants[0].Identity.String(), "Info should contain participant details")
	}

	// Test by non-owner
	infoNonOwner := round.GetParticipantsInfoReplyMessage(nonOwnerID)
	assert.Empty(infoNonOwner, "Expected empty string for non-owner")
}

func TestRound_IsRegistrationClose(t *testing.T) {
	round := NewRound("owner123", "testInvite")
	round.SetIdentity("owner123", Villager, 1)
	assert := assert.New(t)

	assert.False(round.IsRegistrationClose(), "Expected registration to be open initially")

	round.Register("user1", "User One", "url1")
	assert.True(round.IsRegistrationClose(), "Expected registration to be closed after filling all spots")
}

func TestRound_IsRegistrationDuplicate(t *testing.T) {
	round := NewRound("owner123", "testInvite")
	round.SetIdentity("owner123", Villager, 1)
	round.Register("user1", "User One", "url1")
	assert := assert.New(t)

	isDup, p := round.IsRegistrationDuplicate("user1")
	assert.True(isDup, "Expected user1 to be a duplicate")
	assert.NotNil(p, "Expected participant object for user1")
	if p != nil {
		assert.Equal("user1", p.UserID, "Expected participant UID to be user1")
	}

	isDup, p = round.IsRegistrationDuplicate("user2")
	assert.False(isDup, "Expected user2 not to be a duplicate")
	assert.Nil(p, "Expected nil participant object for user2")
}

func TestRound_IsOwner(t *testing.T) {
	ownerID := "owner123"
	round := NewRound(ownerID, "testInvite")
	assert := assert.New(t)

	assert.True(round.IsOwner(ownerID), fmt.Sprintf("Expected %s to be owner", ownerID))
	assert.False(round.IsOwner("nonOwner456"), "Expected nonOwner456 not to be owner")
}

func TestRound_IsExpired(t *testing.T) {
	round := NewRound("owner123", "testInvite")
	assert := assert.New(t)

	// Test not expired
	assert.False(round.IsExpired(), "Expected round not to be expired initially")

	// Test expired
	round.ExpiredAt = time.Now().Add(-1 * time.Hour) // Set expire time to 1 hour ago
	assert.True(round.IsExpired(), "Expected round to be expired")
}
