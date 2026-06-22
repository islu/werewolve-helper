package usecase

import (
	"sync"
	"testing"
	"werewolve-helper/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestRoundManager_PutGetDelete(t *testing.T) {
	m := NewRoundManager()
	assert := assert.New(t)

	_, ok := m.Get("owner1")
	assert.False(ok, "Expected no round before Put")

	round := domain.NewRound("owner1", "123456")
	m.Put(round)

	got, ok := m.Get("owner1")
	assert.True(ok, "Expected round after Put")
	assert.Same(round, got, "Expected the same round instance")

	m.Delete("owner1")
	_, ok = m.Get("owner1")
	assert.False(ok, "Expected no round after Delete")
}

func TestRoundManager_FindByInviteNo(t *testing.T) {
	m := NewRoundManager()
	assert := assert.New(t)

	m.Put(domain.NewRound("owner1", "111111"))
	m.Put(domain.NewRound("owner2", "222222"))

	got, ok := m.FindByInviteNo("222222")
	assert.True(ok, "Expected to find round by invite number")
	assert.Equal("owner2", got.OwnerID, "Expected the round owned by owner2")

	_, ok = m.FindByInviteNo("999999")
	assert.False(ok, "Expected no round for unknown invite number")

	assert.True(m.IsInviteNoDuplicate("111111"), "Expected 111111 to be a duplicate")
	assert.False(m.IsInviteNoDuplicate("999999"), "Expected 999999 not to be a duplicate")
}

func TestRoundManager_ConcurrentAccess(t *testing.T) {
	m := NewRoundManager()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := "owner" + string(rune('A'+n%26))
			m.Put(domain.NewRound(id, "000000"))
			m.Get(id)
			m.FindByInviteNo("000000")
			m.Delete(id)
		}(i)
	}
	wg.Wait()
}
