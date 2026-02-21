// Package game holds the live multiplayer game state (board, players, chat)
// that internal/httpserver will expose over HTTP/WebSocket.
package game

import "github.com/jamesjohnsdev/go-chess/engine"

type ChatMessage struct {
	From string
	Text string
}

type Session struct {
	ID    string
	Board *engine.Board
	White string
	Black string
	Chat  []ChatMessage
}

func NewSession(id, white, black string) *Session {
	return &Session{
		ID:    id,
		Board: engine.NewBoard(),
		White: white,
		Black: black,
	}
}

func (s *Session) AddChat(from, text string) {
	s.Chat = append(s.Chat, ChatMessage{From: from, Text: text})
}
