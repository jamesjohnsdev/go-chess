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

func TestMakeMoveRejectedAfterCheckmate(t *testing.T) {
	s := newSession("test")
	moves := []struct {
		color engine.Color
		uci   string
	}{
		{engine.White, "f2f3"},
		{engine.Black, "e7e5"},
		{engine.White, "g2g4"},
		{engine.Black, "d8h4"},
	}
	for _, mv := range moves {
		m, err := engine.ParseMove(mv.uci)
		if err != nil {
			t.Fatal(err)
		}
		if err := s.MakeMove(mv.color, m); err != nil {
			t.Fatalf("MakeMove(%s): %v", mv.uci, err)
		}
	}

	if !s.State().Checkmate {
		t.Fatal("expected checkmate after fool's mate")
	}

	m, _ := engine.ParseMove("g1f3")
	if err := s.MakeMove(engine.White, m); !errors.Is(err, ErrGameOver) {
		t.Errorf("MakeMove after checkmate = %v, want ErrGameOver", err)
	}
}

func TestMakeMoveRejectedAfterThreefoldRepetition(t *testing.T) {
	s := newSession("test")
	shuffle := []struct {
		color engine.Color
		uci   string
	}{
		{engine.White, "g1f3"},
		{engine.Black, "g8f6"},
		{engine.White, "f3g1"},
		{engine.Black, "f6g8"},
	}
	for rep := 0; rep < 2; rep++ {
		for _, mv := range shuffle {
			m, err := engine.ParseMove(mv.uci)
			if err != nil {
				t.Fatal(err)
			}
			if err := s.MakeMove(mv.color, m); err != nil {
				t.Fatalf("MakeMove(%s): %v", mv.uci, err)
			}
		}
	}

	state := s.State()
	if !state.Draw || state.DrawReason != "threefold repetition" {
		t.Fatalf("State() = %+v, want draw by threefold repetition", state)
	}

	m, _ := engine.ParseMove("b1c3")
	if err := s.MakeMove(engine.White, m); !errors.Is(err, ErrGameOver) {
		t.Errorf("MakeMove after threefold repetition = %v, want ErrGameOver", err)
	}
}

func TestComputerBlackAutoRepliesToHumanMove(t *testing.T) {
	s := newComputerSession("test", engine.Black)

	if got := s.State().Turn; got != "white" {
		t.Fatalf("Turn() right after creation = %q, want white (computer shouldn't move first)", got)
	}

	m, _ := engine.ParseMove("e2e4")
	if err := s.MakeMove(engine.White, m); err != nil {
		t.Fatalf("MakeMove(e2e4): %v", err)
	}

	state := s.State()
	if state.Turn != "white" {
		t.Errorf("Turn() after computer's reply = %q, want white", state.Turn)
	}
	if state.Board == engine.NewBoard().String() {
		t.Error("board unchanged after human + computer moves")
	}
}

func TestComputerWhiteMovesFirst(t *testing.T) {
	s := newComputerSession("test", engine.White)

	state := s.State()
	if state.Turn != "black" {
		t.Fatalf("Turn() right after creation = %q, want black (computer should have moved)", state.Turn)
	}
	if state.Board == engine.NewBoard().String() {
		t.Error("board unchanged; computer should have made the opening move")
	}
}

func TestMakeMoveRejectsComputersColor(t *testing.T) {
	s := newComputerSession("test", engine.Black)
	m, _ := engine.ParseMove("e7e5")
	if err := s.MakeMove(engine.Black, m); !errors.Is(err, ErrNotYourTurn) {
		t.Errorf("MakeMove as the computer's color = %v, want ErrNotYourTurn", err)
	}
}

func TestCreateVsComputerOnlyIssuesHumanToken(t *testing.T) {
	store := NewStore()
	s := store.CreateVsComputer(engine.Black)
	white, black := s.Tokens()
	if white == "" {
		t.Error("human (white) token should be set")
	}
	if black == "" {
		t.Error("computer (black) token should still be generated internally, just never handed out by the API layer")
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
