package ai

import (
	"testing"

	"github.com/jamesjohnsdev/go-chess/engine"
)

func TestEvaluateStartingPositionIsSymmetric(t *testing.T) {
	if got := evaluate(engine.NewBoard()); got != 0 {
		t.Errorf("evaluate(NewBoard()) = %d, want 0 (material and placement are symmetric)", got)
	}
}

func TestPSTKnightPrefersCenterOverCorner(t *testing.T) {
	center := engine.NewEmptyBoard()
	center.Set(sq(t, "e1"), engine.Piece{Type: engine.King, Color: engine.White})
	center.Set(sq(t, "e8"), engine.Piece{Type: engine.King, Color: engine.Black})
	center.Set(sq(t, "e4"), engine.Piece{Type: engine.Knight, Color: engine.White})

	corner := engine.NewEmptyBoard()
	corner.Set(sq(t, "e1"), engine.Piece{Type: engine.King, Color: engine.White})
	corner.Set(sq(t, "e8"), engine.Piece{Type: engine.King, Color: engine.Black})
	corner.Set(sq(t, "a1"), engine.Piece{Type: engine.Knight, Color: engine.White})

	if evaluate(center) <= evaluate(corner) {
		t.Errorf("evaluate(knight on e4) = %d, want > evaluate(knight on a1) = %d", evaluate(center), evaluate(corner))
	}
}

func TestPSTAdvancedCentralPawnScoresHigher(t *testing.T) {
	advanced := engine.NewEmptyBoard()
	advanced.Set(sq(t, "e1"), engine.Piece{Type: engine.King, Color: engine.White})
	advanced.Set(sq(t, "e8"), engine.Piece{Type: engine.King, Color: engine.Black})
	advanced.Set(sq(t, "e5"), engine.Piece{Type: engine.Pawn, Color: engine.White})

	start := engine.NewEmptyBoard()
	start.Set(sq(t, "e1"), engine.Piece{Type: engine.King, Color: engine.White})
	start.Set(sq(t, "e8"), engine.Piece{Type: engine.King, Color: engine.Black})
	start.Set(sq(t, "e2"), engine.Piece{Type: engine.Pawn, Color: engine.White})

	if evaluate(advanced) <= evaluate(start) {
		t.Errorf("evaluate(pawn on e5) = %d, want > evaluate(pawn on e2) = %d", evaluate(advanced), evaluate(start))
	}
}

func TestPSTMirrorsForBlack(t *testing.T) {
	// e7 is Black's mirror of White's e2: same distance from its own back
	// rank, so their piece-square bonuses should match exactly.
	white := pstValue(engine.Piece{Type: engine.Pawn, Color: engine.White}, sq(t, "e2"))
	black := pstValue(engine.Piece{Type: engine.Pawn, Color: engine.Black}, sq(t, "e7"))
	if white != black {
		t.Errorf("pstValue(white pawn, e2) = %d, pstValue(black pawn, e7) = %d, want equal", white, black)
	}
}
