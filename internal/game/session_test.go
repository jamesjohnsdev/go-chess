package game

import (
	"errors"
	"testing"
	"time"

	"github.com/jamesjohnsdev/go-chess/engine"
)

func TestColorForToken(t *testing.T) {
	s := newSession("test")
	white, black := s.Tokens()

	if c, err := s.ColorForToken(white); err != nil || c != engine.White {
		t.Errorf("ColorForToken(white) = %v, %v, want White, nil", c, err)
	}
	if c, err := s.ColorForToken(black); err != nil || c != engine.Black {
		t.Errorf("ColorForToken(black) = %v, %v, want Black, nil", c, err)
	}
	if _, err := s.ColorForToken("bogus"); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("ColorForToken(bogus) err = %v, want ErrInvalidToken", err)
	}
}

func TestMakeMoveWrongTurn(t *testing.T) {
	s := newSession("test")
	m, _ := engine.ParseMove("e7e5")
	if err := s.MakeMove(engine.Black, m); !errors.Is(err, ErrNotYourTurn) {
		t.Errorf("MakeMove out of turn = %v, want ErrNotYourTurn", err)
	}
}

func TestMakeMoveIllegal(t *testing.T) {
	s := newSession("test")
	m, _ := engine.ParseMove("e2e5")
	if err := s.MakeMove(engine.White, m); !errors.Is(err, ErrIllegalMove) {
		t.Errorf("MakeMove(e2e5) = %v, want ErrIllegalMove", err)
	}
}

func TestMakeMoveLegalUpdatesState(t *testing.T) {
	s := newSession("test")
	m, _ := engine.ParseMove("e2e4")
	if err := s.MakeMove(engine.White, m); err != nil {
		t.Fatalf("MakeMove(e2e4) = %v, want nil", err)
	}
	if got := s.State().Turn; got != "black" {
		t.Errorf("Turn() after e2e4 = %q, want black", got)
	}
}

func TestSubscribeReceivesEvents(t *testing.T) {
	s := newSession("test")
	events, cancel := s.Subscribe()
	defer cancel()

	select {
	case e := <-events:
		if e.Type != "state" {
			t.Errorf("first event type = %q, want state", e.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for initial state event")
	}

	m, _ := engine.ParseMove("e2e4")
	if err := s.MakeMove(engine.White, m); err != nil {
		t.Fatalf("MakeMove: %v", err)
	}

	select {
	case e := <-events:
		if e.Type != "move" || e.Move == nil || *e.Move != (MoveMsg{From: "e2", To: "e4"}) {
			t.Errorf("move event = %+v, want move e2e4", e)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for move event")
	}
}

func TestAddChatBroadcasts(t *testing.T) {
	s := newSession("test")
	events, cancel := s.Subscribe()
	defer cancel()
	<-events // initial state

	s.AddChat("white", "gg")

	select {
	case e := <-events:
		if e.Type != "chat" || e.Chat == nil || *e.Chat != (ChatMessage{From: "white", Text: "gg"}) {
			t.Errorf("chat event = %+v, want chat from white", e)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for chat event")
	}
	if len(s.State().Chat) != 1 {
		t.Errorf("State().Chat = %v, want 1 message", s.State().Chat)
	}
}

func TestCancelUnsubscribes(t *testing.T) {
	s := newSession("test")
	events, cancel := s.Subscribe()
	<-events // initial state
	cancel()

	if _, open := <-events; open {
		t.Error("channel should be closed after cancel")
	}

	m, _ := engine.ParseMove("e2e4")
	if err := s.MakeMove(engine.White, m); err != nil {
		t.Fatalf("MakeMove: %v", err)
	}
}
