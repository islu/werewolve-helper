package internal

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Round struct {
	// ID               int64
	// UID              string
	OwnerID          string
	InviteNo         string
	CreateTime       time.Time
	ExpireTime       time.Time
	Participants     []Participant
	Identities       []Identity
	TempIdentity     Identity
	TempIdentityFlag bool
}

func NewRound(userID, inviteNo string) *Round {
	return &Round{
		OwnerID:          userID,
		InviteNo:         inviteNo,
		Identities:       []Identity{},
		Participants:     []Participant{},
		CreateTime:       time.Now(),
		ExpireTime:       time.Now().Add(2 * time.Hour),
		TempIdentityFlag: false,
	}
}

func NewRoundWith9PersonStandardMode(userID, inviteNo string) *Round {

	round := NewRound(userID, inviteNo)

	round.SetIdentity(userID, Seer, 1)
	round.SetIdentity(userID, Witch, 1)
	round.SetIdentity(userID, Hunter, 1)
	round.SetIdentity(userID, Villager, 3)
	round.SetIdentity(userID, Werewolf, 3)

	return round
}

func (r *Round) SetIdentity(userID string, iden Identity, nums int) {
	if r.IsOwner(userID) {
		for i := 0; i < nums; i++ {
			r.Identities = append(r.Identities, iden)
		}
		rng.Shuffle(len(r.Identities), func(i, j int) {
			r.Identities[i], r.Identities[j] = r.Identities[j], r.Identities[i]
		})
	}
}

func (r *Round) Register(userID, name, pictureURL string) string {
	if r.IsRegistrationClose() {
		log.Println("register is closed")
		return ""
	}

	idx := len(r.Participants)
	user := NewParticipant(userID, name, pictureURL, r.Identities[idx])
	r.Participants = append(r.Participants, *user)
	return r.Identities[idx].String()
}

func (r *Round) Again() {
	rng.Shuffle(len(r.Identities), func(i, j int) {
		r.Identities[i], r.Identities[j] = r.Identities[j], r.Identities[i]
	})
	for i, p := range r.Participants {
		p.Identity = r.Identities[i]
	}
}

func (r *Round) GetParticipantsInfoStr(userID string) string {
	if !r.IsOwner(userID) {
		log.Println(r.InviteNo)
		log.Println(r.OwnerID + " not equal " + userID)
		return ""
	}

	var sb strings.Builder
	sb.WriteString("目前參與人數: ")
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

func (r *Round) IsRegistrationClose() bool {
	return len(r.Identities) == len(r.Participants)
}

func (r *Round) IsOwner(creatorID string) bool {
	return r.OwnerID == creatorID
}

func (r *Round) IsExpired() bool {
	return r.ExpireTime.Before(time.Now())
}

type Participant struct {
	// ID         int64
	UID        string
	Name       string
	PictureURL string
	Identity   Identity
}

func NewParticipant(userID, name, pictureURL string, identity Identity) *Participant {
	return &Participant{
		UID:        userID,
		Name:       name,
		PictureURL: pictureURL,
		Identity:   identity,
	}
}

type Identity int

const (
	Werewolf Identity = iota + 1
	Villager
	Seer
	Witch
	Hunter
)

func (iden Identity) String() string {
	switch iden {
	case Werewolf:
		return "狼人"
	case Villager:
		return "平民"
	case Seer:
		return "預言家"
	case Witch:
		return "女巫"
	case Hunter:
		return "獵人"
	}
	return "unknown"
}
