package domain

// Identity represents the role of a player in the game.
type Identity int

// Constants for different identities (roles) in the game.
const (
	WerewolfKing   Identity = iota + 1 // Werewolf King, Wolf team
	WhiteWerewolf                      // White Werewolf, Wolf team
	GhostRider                         // Ghost Rider, Wolf team
	WerewolfBeauty                     // Werewolf Beauty, Wolf team
	Werewolf                           // Werewolf, Wolf team
	Seer                               // Seer, Villager team
	Witch                              // Witch, Villager team
	Hunter                             // Hunter, Villager team
	Guard                              // Guard, Villager team
	Knight                             // Knight, Villager team
	Magician                           // Magician, Villager team
	Villager                           // Villager, Villager team
)

// String returns the string representation of an Identity.
func (iden Identity) String() string {
	switch iden {
	case WerewolfKing:
		return "狼王"
	case WhiteWerewolf:
		return "白狼王"
	case GhostRider:
		return "惡靈騎士"
	case WerewolfBeauty:
		return "狼美人"
	case Werewolf:
		return "狼人"
	case Seer:
		return "預言家"
	case Witch:
		return "女巫"
	case Hunter:
		return "獵人"
	case Guard:
		return "守衛"
	case Knight:
		return "騎士"
	case Magician:
		return "魔術師"
	case Villager:
		return "平民"
	default: // Default for unhandled identities.
		return "unknown"
	}
}

// Participant represents a player in the game.
type Participant struct {
	UserID     string   // User ID of the participant.
	Name       string   // Name of the participant.
	PictureURL string   // URL of the participant's picture.
	Identity   Identity // Assigned identity (role) of the participant.
}

// NewParticipant creates a new participant.
func NewParticipant(userID, name, pictureURL string, identity Identity) *Participant {
	return &Participant{
		UserID:     userID,
		Name:       name,
		PictureURL: pictureURL,
		Identity:   identity,
	}
}
