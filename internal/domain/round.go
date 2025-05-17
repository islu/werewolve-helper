package domain

import (
	"crypto/rand"
	"encoding/binary"
	"log"
	mrand "math/rand/v2" // Using math/rand/v2
	"strconv"
	"strings"
	"time"
)

var (
	// rng is a random number generator used for shuffling, seeded with crypto/rand.
	rng *mrand.Rand
)

func init() {
	// Seed math/rand/v2 with a cryptographically secure random number.
	var seed uint64 // math/rand/v2 NewSource expects a uint64 seed
	err := binary.Read(rand.Reader, binary.BigEndian, &seed)
	if err != nil {
		// Fallback to time-based seed if crypto/rand fails, though this is less ideal.
		log.Printf("Failed to read crypto/rand for seed, falling back to time-based seed: %v", err)
		seed = uint64(time.Now().UnixNano())
	}
	rng = mrand.New(mrand.NewPCG(seed, seed>>32)) // Example using PCG source, or use NewChaCha8
	// or simply randv2.New(randv2.NewSource(seed)) if a default Source is acceptable
	// For direct equivalent of old mrand.NewSource(seed):
	// rng = randv2.New(randv2.NewSource(seed))
	// However, math/rand/v2 encourages using specific sources like PCG or ChaCha8.
	// Using PCG as an example:
	source := mrand.NewPCG(seed, seed+1) // PCG requires two uint64 seeds
	rng = mrand.New(source)

	// If you want to keep it simple and closer to the old math/rand behavior with a single seed:
	// var simpleSeed uint64
	// err_simple := binary.Read(crypto_rand.Reader, binary.BigEndian, &simpleSeed)
	// if err_simple != nil {
	//  log.Printf("Failed to read crypto/rand for simple seed, falling back to time-based seed: %v", err_simple)
	//  simpleSeed = uint64(time.Now().UnixNano())
	// }
	// rng = randv2.New(randv2.NewSource(simpleSeed)) // This uses the default source type for NewSource
}

// Round represents a game round.
type Round struct {
	OwnerID          string        // ID of the user who created the round.
	InviteNo         string        // Invitation number for the round.
	CreatedAt        time.Time     // Time when the round was created.
	ExpiredAt        time.Time     // Time when the round expires.
	Participants     []Participant // List of participants in the round.
	Identities       []Identity    // List of identities (roles) assigned in the round.
	TempIdentity     Identity
	TempIdentityFlag bool
}

// NewRound creates a new game round.
func NewRound(userID, inviteNo string) *Round {
	return &Round{
		OwnerID:          userID,
		InviteNo:         inviteNo,
		Identities:       []Identity{},
		Participants:     []Participant{},
		CreatedAt:        time.Now(),
		ExpiredAt:        time.Now().Add(2 * time.Hour), // Round expires in 2 hours.
		TempIdentityFlag: false,
	}
}

// SetIdentity sets the identities for the round. Only the owner can set identities.
// It adds the specified identity 'nums' times and shuffles the list of identities.
func (r *Round) SetIdentity(userID string, iden Identity, nums int) {
	if r.IsOwner(userID) {
		for range nums {
			r.Identities = append(r.Identities, iden)
		}
		// Shuffle identities to randomize assignment.
		rng.Shuffle(len(r.Identities), func(i, j int) {
			r.Identities[i], r.Identities[j] = r.Identities[j], r.Identities[i]
		})
	}
}

// Register allows a user to join the round.
// If registration is closed or the user is already registered, it returns an empty string.
// Otherwise, it assigns an identity to the user and returns the identity as a string.
func (r *Round) Register(userID, name, pictureURL string) string {
	if r.IsRegistrationClose() {
		log.Println("register is closed")
		return ""
	}

	// Assign the next available identity.
	idx := len(r.Participants)
	user := NewParticipant(userID, name, pictureURL, r.Identities[idx])
	r.Participants = append(r.Participants, *user)
	return r.Identities[idx].String()
}

// Again resets the round for a new game with the same identities.
// It shuffles identities, clears participants, and extends the expiration time.
func (r *Round) Again() {
	rng.Shuffle(len(r.Identities), func(i, j int) {
		r.Identities[i], r.Identities[j] = r.Identities[j], r.Identities[i]
	})
	// Empty participants for the new game.
	r.Participants = []Participant{}
	// Extend expire time for the new game.
	r.ExpiredAt = time.Now().Add(2 * time.Hour)
}

// GetParticipantsInfoReplyMessage returns a string with information about participants.
// Only the owner can get this information.
// The string includes the count of participants and their assigned identities.
func (r *Round) GetParticipantsInfoReplyMessage(userID string) string {

	// Check if the user is the owner of the round.
	if !r.IsOwner(userID) {
		log.Println(r.InviteNo)
		log.Println(r.OwnerID + " not equal " + userID)
		return ""
	}

	var sb strings.Builder
	sb.WriteString("目前參與人數: ") // "Current number of participants: "
	sb.WriteString(strconv.Itoa(len(r.Participants)))
	sb.WriteString("/")
	sb.WriteString(strconv.Itoa(len(r.Identities)))

	for _, p := range r.Participants {
		sb.WriteString("\n")
		sb.WriteString(p.Name)
		sb.WriteString(":")
		sb.WriteString(p.Identity.String())
	}
	return sb.String()
}

//
// --- Validation functions ---
//

// IsRegistrationClose checks if the registration for the round is closed.
// Registration is closed when the number of participants equals the number of identities.
func (r *Round) IsRegistrationClose() bool {
	return len(r.Identities) == len(r.Participants)
}

// IsRegistrationDuplicate checks if a user is already registered in the round.
// It returns true and the participant if duplicate, false and nil otherwise.
func (r *Round) IsRegistrationDuplicate(userID string) (bool, *Participant) {
	for _, p := range r.Participants {
		if p.UserID == userID {
			return true, &p
		}
	}
	return false, nil
}

// IsOwner checks if the given userID is the owner of the round.
func (r *Round) IsOwner(creatorID string) bool {
	return r.OwnerID == creatorID
}

// IsExpired checks if the round has expired.
func (r *Round) IsExpired() bool {
	return r.ExpiredAt.Before(time.Now())
}
