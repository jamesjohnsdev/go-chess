package game

import (
	"sync"

	"github.com/jamesjohnsdev/go-chess/engine"
)

// Store holds every in-flight Session, keyed by ID. It's an in-memory store:
// games are lost on restart, which is fine until go-chess needs persistence.
type Store struct {
	mu       sync.Mutex
	sessions map[string]*Session
}

func NewStore() *Store {
	return &Store{sessions: make(map[string]*Session)}
}

func (st *Store) Create() *Session {
	st.mu.Lock()
	defer st.mu.Unlock()

	s := newSession(randomID(8))
	st.sessions[s.ID] = s
	return s
}

// CreateVsComputer creates a game where computerColor is played by the
// built-in AI opponent.
func (st *Store) CreateVsComputer(computerColor engine.Color) *Session {
	st.mu.Lock()
	defer st.mu.Unlock()

	s := newComputerSession(randomID(8), computerColor)
	st.sessions[s.ID] = s
	return s
}

func (st *Store) Get(id string) (*Session, bool) {
	st.mu.Lock()
	defer st.mu.Unlock()

	s, ok := st.sessions[id]
	return s, ok
}
