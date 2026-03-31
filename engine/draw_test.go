package engine

import "testing"

func TestFiftyMoveRule(t *testing.T) {
	b := NewEmptyBoard()
	b.Set(mustSquare(t, "a1"), Piece{Type: King, Color: White})
	b.Set(mustSquare(t, "h8"), Piece{Type: King, Color: Black})
	b.Set(mustSquare(t, "a2"), Piece{Type: Pawn, Color: White})

	shuffle := []Move{
		{From: mustSquare(t, "a1"), To: mustSquare(t, "b1")},
		{From: mustSquare(t, "h8"), To: mustSquare(t, "h7")},
		{From: mustSquare(t, "b1"), To: mustSquare(t, "a1")},
		{From: mustSquare(t, "h7"), To: mustSquare(t, "h8")},
	}
	for i := 0; i < 99; i++ {
		m := shuffle[i%len(shuffle)]
		b.MakeMove(m)
		if b.IsFiftyMoveDraw() {
			t.Fatalf("IsFiftyMoveDraw() = true after %d halfmoves, too early", i+1)
		}
	}
	b.MakeMove(shuffle[99%len(shuffle)])
	if !b.IsFiftyMoveDraw() {
		t.Errorf("IsFiftyMoveDraw() = false after 100 halfmoves, want true")
	}

	// A pawn move resets the clock.
	b2 := NewBoard()
	b2.MakeMove(Move{From: mustSquare(t, "e2"), To: mustSquare(t, "e4")})
	if got := b2.HalfmoveClock(); got != 0 {
		t.Errorf("HalfmoveClock() after pawn move = %d, want 0", got)
	}
}

func TestInsufficientMaterial(t *testing.T) {
	tests := []struct {
		name  string
		setup func(b *Board)
		want  bool
	}{
		{"bare kings", func(b *Board) {}, true},
		{"king and knight vs king", func(b *Board) {
			b.Set(mustSquare(t, "b1"), Piece{Type: Knight, Color: White})
		}, true},
		{"king and bishop vs king", func(b *Board) {
			b.Set(mustSquare(t, "c1"), Piece{Type: Bishop, Color: White})
		}, true},
		{"same-color bishops", func(b *Board) {
			b.Set(mustSquare(t, "c1"), Piece{Type: Bishop, Color: White})
			b.Set(mustSquare(t, "f8"), Piece{Type: Bishop, Color: Black})
		}, true},
		{"opposite-color bishops", func(b *Board) {
			b.Set(mustSquare(t, "c1"), Piece{Type: Bishop, Color: White})
			b.Set(mustSquare(t, "c8"), Piece{Type: Bishop, Color: Black})
		}, false},
		{"king and rook vs king", func(b *Board) {
			b.Set(mustSquare(t, "a1"), Piece{Type: Rook, Color: White})
		}, false},
		{"two knights", func(b *Board) {
			b.Set(mustSquare(t, "b1"), Piece{Type: Knight, Color: White})
			b.Set(mustSquare(t, "g1"), Piece{Type: Knight, Color: White})
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewEmptyBoard()
			b.Set(mustSquare(t, "e1"), Piece{Type: King, Color: White})
			b.Set(mustSquare(t, "e8"), Piece{Type: King, Color: Black})
			tt.setup(b)
			if got := b.IsInsufficientMaterial(); got != tt.want {
				t.Errorf("IsInsufficientMaterial() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPositionKeyRepetition(t *testing.T) {
	b := NewBoard()
	counts := map[string]int{b.PositionKey(): 1}

	shuffle := []string{"g1f3", "g8f6", "f3g1", "f6g8"}
	for rep := 0; rep < 2; rep++ {
		for _, uci := range shuffle {
			m, err := ParseMove(uci)
			if err != nil {
				t.Fatal(err)
			}
			b.MakeMove(m)
			counts[b.PositionKey()]++
		}
	}

	if got := counts[b.PositionKey()]; got != 3 {
		t.Errorf("position occurred %d times, want 3 (threefold repetition)", got)
	}
}

func TestPositionKeyDistinguishesCastlingRights(t *testing.T) {
	a := NewEmptyBoard()
	a.Set(mustSquare(t, "e1"), Piece{Type: King, Color: White})
	a.Set(mustSquare(t, "h1"), Piece{Type: Rook, Color: White})
	a.castling.WhiteKingside = true

	b := NewEmptyBoard()
	b.Set(mustSquare(t, "e1"), Piece{Type: King, Color: White})
	b.Set(mustSquare(t, "h1"), Piece{Type: Rook, Color: White})
	b.castling.WhiteKingside = false

	if a.PositionKey() == b.PositionKey() {
		t.Error("PositionKey() should differ when castling rights differ")
	}
}
