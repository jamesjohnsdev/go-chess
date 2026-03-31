// Package game holds live multiplayer game state (board, players, chat) and
// a pub-sub Event stream that internal/httpserver exposes over WebSocket,
// plus REST endpoints for creating and inspecting games.
package game

import (
	"errors"
	"sync"

	"github.com/jamesjohnsdev/go-chess/engine"
)

var (
	ErrNotYourTurn  = errors.New("not your turn")
	ErrIllegalMove  = errors.New("illegal move")
	ErrInvalidToken = errors.New("invalid token")
	ErrGameOver     = errors.New("game is already over")
)

type ChatMessage struct {
	From string `json:"from"`
	Text string `json:"text"`
}

// MoveMsg is the wire representation of a move: algebraic squares plus an
// optional lowercase promotion letter, e.g. {From: "e7", To: "e8", Promotion: "q"}.
type MoveMsg struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Promotion string `json:"promotion,omitempty"`
}

func (m MoveMsg) ToMove() (engine.Move, error) {
	return engine.ParseMove(m.From + m.To + m.Promotion)
}

func moveToMsg(m engine.Move) MoveMsg {
	return MoveMsg{From: m.From.String(), To: m.To.String(), Promotion: m.Promotion.Letter()}
}

// ClientMessage is a message sent by a connected player over WebSocket.
type ClientMessage struct {
	Type string   `json:"type"` // "move" or "chat"
	Move *MoveMsg `json:"move,omitempty"`
	Text string   `json:"text,omitempty"`
}

// State is a snapshot of a game, safe to serialize and send to clients.
type State struct {
	Turn       string        `json:"turn"`
	Board      string        `json:"board"`
	Check      bool          `json:"check"`
	Checkmate  bool          `json:"checkmate"`
	Stalemate  bool          `json:"stalemate"`
	Draw       bool          `json:"draw"`
	DrawReason string        `json:"draw_reason,omitempty"`
	Chat       []ChatMessage `json:"chat"`
}

func (s State) Over() bool {
	return s.Checkmate || s.Stalemate || s.Draw
}

// Event is broadcast to every subscriber whenever a game changes.
type Event struct {
	Type  string       `json:"type"` // "state", "move", "chat", or "error"
	Move  *MoveMsg     `json:"move,omitempty"`
	Chat  *ChatMessage `json:"chat,omitempty"`
	Error string       `json:"error,omitempty"`
	State State        `json:"state"`
}

// Session is one live game: board state, the two player tokens, chat
// history, and the set of subscribers to broadcast events to. All access
// goes through its methods, which are safe for concurrent use.
type Session struct {
	ID string

	mu         sync.Mutex
	board      *engine.Board
	whiteToken string
	blackToken string
	chat       []ChatMessage
	subs       map[int]chan Event
	nextSub    int
	positions  map[string]int
}

func newSession(id string) *Session {
	board := engine.NewBoard()
	return &Session{
		ID:         id,
		board:      board,
		whiteToken: randomID(16),
		blackToken: randomID(16),
		subs:       make(map[int]chan Event),
		positions:  map[string]int{board.PositionKey(): 1},
	}
}

// Tokens returns the bearer tokens a client uses to play as white or black.
func (s *Session) Tokens() (white, black string) {
	return s.whiteToken, s.blackToken
}

func (s *Session) ColorForToken(token string) (engine.Color, error) {
	switch token {
	case s.whiteToken:
		return engine.White, nil
	case s.blackToken:
		return engine.Black, nil
	default:
		return engine.White, ErrInvalidToken
	}
}

func (s *Session) State() State {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stateLocked()
}

func (s *Session) stateLocked() State {
	draw, reason := s.drawStatusLocked()
	return State{
		Turn:       s.board.Turn().String(),
		Board:      s.board.String(),
		Check:      s.board.InCheck(),
		Checkmate:  s.board.IsCheckmate(),
		Stalemate:  s.board.IsStalemate(),
		Draw:       draw,
		DrawReason: reason,
		Chat:       append([]ChatMessage(nil), s.chat...),
	}
}

// drawStatusLocked reports draws other than stalemate, which State already
// surfaces via its own field.
func (s *Session) drawStatusLocked() (bool, string) {
	switch {
	case s.board.IsFiftyMoveDraw():
		return true, "fifty-move rule"
	case s.board.IsInsufficientMaterial():
		return true, "insufficient material"
	case s.positions[s.board.PositionKey()] >= 3:
		return true, "threefold repetition"
	default:
		return false, ""
	}
}

// MakeMove applies m as color's move if it's their turn and the move is
// legal, then broadcasts the resulting state to all subscribers.
func (s *Session) MakeMove(color engine.Color, m engine.Move) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.board.IsCheckmate() || s.board.IsStalemate() {
		return ErrGameOver
	}
	if draw, _ := s.drawStatusLocked(); draw {
		return ErrGameOver
	}
	if s.board.Turn() != color {
		return ErrNotYourTurn
	}
	legal := false
	for _, lm := range s.board.LegalMoves() {
		if lm == m {
			legal = true
			break
		}
	}
	if !legal {
		return ErrIllegalMove
	}

	s.board.MakeMove(m)
	s.positions[s.board.PositionKey()]++
	msg := moveToMsg(m)
	s.broadcastLocked(Event{Type: "move", Move: &msg, State: s.stateLocked()})
	return nil
}

// AddChat appends a chat message and broadcasts it to all subscribers.
func (s *Session) AddChat(from, text string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg := ChatMessage{From: from, Text: text}
	s.chat = append(s.chat, msg)
	s.broadcastLocked(Event{Type: "chat", Chat: &msg, State: s.stateLocked()})
}

// Subscribe registers a new event listener, immediately queuing the current
// state. Call the returned cancel func to unsubscribe and close the channel.
func (s *Session) Subscribe() (<-chan Event, func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextSub
	s.nextSub++
	ch := make(chan Event, 16)
	s.subs[id] = ch
	ch <- Event{Type: "state", State: s.stateLocked()}

	cancel := func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if ch, ok := s.subs[id]; ok {
			delete(s.subs, id)
			close(ch)
		}
	}
	return ch, cancel
}

// broadcastLocked drops the event for any subscriber whose channel is full
// rather than blocking the whole session on a slow reader.
func (s *Session) broadcastLocked(e Event) {
	for _, ch := range s.subs {
		select {
		case ch <- e:
		default:
		}
	}
}
