package usecase

import (
	"sync"
	"werewolve-helper/internal/domain"
)

// RoundManager is a concurrency-safe store for game rounds, keyed by owner ID.
// LINE webhook events are handled in separate goroutines, so all access to the
// underlying map is guarded by a mutex.
type RoundManager struct {
	mu     sync.RWMutex
	rounds map[string]*domain.Round
}

// NewRoundManager creates an empty RoundManager.
func NewRoundManager() *RoundManager {
	return &RoundManager{
		rounds: make(map[string]*domain.Round),
	}
}

// Get returns the round owned by ownerID, if one exists.
func (m *RoundManager) Get(ownerID string) (*domain.Round, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.rounds[ownerID]
	return r, ok
}

// Put stores a round, keyed by its owner ID.
func (m *RoundManager) Put(r *domain.Round) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rounds[r.OwnerID] = r
}

// Delete removes the round owned by ownerID.
func (m *RoundManager) Delete(ownerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.rounds, ownerID)
}

// FindByInviteNo returns the round whose invite number matches inviteNo.
func (m *RoundManager) FindByInviteNo(inviteNo string) (*domain.Round, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, r := range m.rounds {
		if r.InviteNo == inviteNo {
			return r, true
		}
	}
	return nil, false
}

// IsInviteNoDuplicate reports whether any round already uses inviteNo.
func (m *RoundManager) IsInviteNoDuplicate(inviteNo string) bool {
	_, ok := m.FindByInviteNo(inviteNo)
	return ok
}
