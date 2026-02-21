package engine

import "testing"

func TestNewBoardStartPosition(t *testing.T) {
	b := NewBoard()

	if b.Turn() != White {
		t.Errorf("Turn() = %v, want White", b.Turn())
	}

	tests := []struct {
		sq   string
		want Piece
	}{
		{"a1", Piece{Type: Rook, Color: White}},
		{"e1", Piece{Type: King, Color: White}},
		{"e2", Piece{Type: Pawn, Color: White}},
		{"e8", Piece{Type: King, Color: Black}},
		{"h8", Piece{Type: Rook, Color: Black}},
		{"e4", Piece{}},
	}
	for _, tt := range tests {
		sq, err := ParseSquare(tt.sq)
		if err != nil {
			t.Fatalf("ParseSquare(%q): %v", tt.sq, err)
		}
		if got := b.At(sq); got != tt.want {
			t.Errorf("At(%s) = %+v, want %+v", tt.sq, got, tt.want)
		}
	}
}

func TestMakeMove(t *testing.T) {
	b := NewBoard()
	from, _ := ParseSquare("e2")
	to, _ := ParseSquare("e4")

	b.MakeMove(Move{From: from, To: to})

	if got := b.At(from); got != (Piece{}) {
		t.Errorf("At(e2) after move = %+v, want empty", got)
	}
	if got := b.At(to); got != (Piece{Type: Pawn, Color: White}) {
		t.Errorf("At(e4) after move = %+v, want white pawn", got)
	}
	if b.Turn() != Black {
		t.Errorf("Turn() after move = %v, want Black", b.Turn())
	}
}
