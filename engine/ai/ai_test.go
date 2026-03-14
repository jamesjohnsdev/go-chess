package ai

import (
	"testing"

	"github.com/jamesjohnsdev/go-chess/engine"
)

func TestBestMovePlaysLegalMoveFromStart(t *testing.T) {
	b := engine.NewBoard()
	e := New()
	e.Depth = 2

	move, ok := e.BestMove(b)
	if !ok {
		t.Fatal("expected a move from the starting position")
	}
	legal := false
	for _, m := range b.LegalMoves() {
		if m == move {
			legal = true
			break
		}
	}
	if !legal {
		t.Errorf("BestMove() = %v, not in LegalMoves()", move)
	}
}

func TestBestMoveFindsMateInOne(t *testing.T) {
	b := engine.NewEmptyBoard()
	b.Set(sq(t, "h8"), engine.Piece{Type: engine.King, Color: engine.Black})
	b.Set(sq(t, "g6"), engine.Piece{Type: engine.King, Color: engine.White})
	b.Set(sq(t, "a1"), engine.Piece{Type: engine.Queen, Color: engine.White})

	e := New()
	e.Depth = 2

	move, ok := e.BestMove(b)
	if !ok {
		t.Fatal("expected a move")
	}
	clone := *b
	clone.MakeMove(move)
	if !clone.IsCheckmate() {
		t.Errorf("BestMove() = %v, did not deliver checkmate:\n%s", move, &clone)
	}
}

func TestBestMoveNoMovesWhenGameOver(t *testing.T) {
	b := engine.NewEmptyBoard()
	b.Set(sq(t, "h1"), engine.Piece{Type: engine.King, Color: engine.White})
	b.Set(sq(t, "g3"), engine.Piece{Type: engine.King, Color: engine.Black})
	b.Set(sq(t, "g2"), engine.Piece{Type: engine.Queen, Color: engine.Black})

	if _, ok := New().BestMove(b); ok {
		t.Error("expected no move for a side already checkmated")
	}
}

func sq(t *testing.T, s string) engine.Square {
	t.Helper()
	square, err := engine.ParseSquare(s)
	if err != nil {
		t.Fatalf("ParseSquare(%q): %v", s, err)
	}
	return square
}
