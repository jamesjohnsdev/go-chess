package engine

import "testing"

func perft(b *Board, depth int) int {
	if depth == 0 {
		return 1
	}
	moves := b.LegalMoves()
	if depth == 1 {
		return len(moves)
	}
	count := 0
	for _, m := range moves {
		clone := *b
		clone.MakeMove(m)
		count += perft(&clone, depth-1)
	}
	return count
}

func TestPerftStartingPosition(t *testing.T) {
	// Reference values: https://www.chessprogramming.org/Perft_Results
	want := []int{1, 20, 400, 8902, 197281}
	b := NewBoard()
	for depth, w := range want {
		if got := perft(b, depth); got != w {
			t.Errorf("perft(%d) = %d, want %d", depth, got, w)
		}
	}
}

func TestPromotion(t *testing.T) {
	b := NewEmptyBoard()
	b.Set(mustSquare(t, "e7"), Piece{Type: Pawn, Color: White})
	b.Set(mustSquare(t, "e1"), Piece{Type: King, Color: White})
	b.Set(mustSquare(t, "a8"), Piece{Type: King, Color: Black})

	moves := b.LegalMoves()
	promos := map[PieceType]bool{}
	for _, m := range moves {
		if m.From == mustSquare(t, "e7") && m.To == mustSquare(t, "e8") {
			promos[m.Promotion] = true
		}
	}
	for _, want := range []PieceType{Queen, Rook, Bishop, Knight} {
		if !promos[want] {
			t.Errorf("missing promotion to %v", want)
		}
	}
}

func TestEnPassant(t *testing.T) {
	b := NewEmptyBoard()
	b.Set(mustSquare(t, "e5"), Piece{Type: Pawn, Color: White})
	b.Set(mustSquare(t, "d7"), Piece{Type: Pawn, Color: Black})
	b.Set(mustSquare(t, "e1"), Piece{Type: King, Color: White})
	b.Set(mustSquare(t, "a8"), Piece{Type: King, Color: Black})
	b.turn = Black

	b.MakeMove(Move{From: mustSquare(t, "d7"), To: mustSquare(t, "d5")})

	target, ok := b.EnPassant()
	if !ok || target != mustSquare(t, "d6") {
		t.Fatalf("EnPassant() = %v, %v, want d6, true", target, ok)
	}

	b.MakeMove(Move{From: mustSquare(t, "e5"), To: mustSquare(t, "d6")})

	if got := b.At(mustSquare(t, "d5")); got.Type != None {
		t.Errorf("At(d5) after en passant = %+v, want empty", got)
	}
	if got := b.At(mustSquare(t, "d6")); got != (Piece{Type: Pawn, Color: White}) {
		t.Errorf("At(d6) after en passant = %+v, want white pawn", got)
	}
}

func TestCastlingKingside(t *testing.T) {
	b := NewEmptyBoard()
	b.Set(mustSquare(t, "e1"), Piece{Type: King, Color: White})
	b.Set(mustSquare(t, "h1"), Piece{Type: Rook, Color: White})
	b.Set(mustSquare(t, "a8"), Piece{Type: King, Color: Black})
	b.castling.WhiteKingside = true

	b.MakeMove(Move{From: mustSquare(t, "e1"), To: mustSquare(t, "g1")})

	if got := b.At(mustSquare(t, "g1")); got != (Piece{Type: King, Color: White}) {
		t.Errorf("At(g1) = %+v, want white king", got)
	}
	if got := b.At(mustSquare(t, "f1")); got != (Piece{Type: Rook, Color: White}) {
		t.Errorf("At(f1) = %+v, want white rook", got)
	}
	if b.castling.WhiteKingside || b.castling.WhiteQueenside {
		t.Error("castling rights should be cleared after castling")
	}
}

func TestCastlingBlockedThroughCheck(t *testing.T) {
	b := NewEmptyBoard()
	b.Set(mustSquare(t, "e1"), Piece{Type: King, Color: White})
	b.Set(mustSquare(t, "h1"), Piece{Type: Rook, Color: White})
	b.Set(mustSquare(t, "f8"), Piece{Type: Rook, Color: Black})
	b.Set(mustSquare(t, "a8"), Piece{Type: King, Color: Black})
	b.castling.WhiteKingside = true

	for _, m := range b.LegalMoves() {
		if m.From == mustSquare(t, "e1") && m.To == mustSquare(t, "g1") {
			t.Fatal("castling through an attacked square should not be legal")
		}
	}
}

func TestStalemate(t *testing.T) {
	b := NewEmptyBoard()
	b.Set(mustSquare(t, "a8"), Piece{Type: King, Color: Black})
	b.Set(mustSquare(t, "b6"), Piece{Type: King, Color: White})
	b.Set(mustSquare(t, "c7"), Piece{Type: Queen, Color: White})
	b.turn = Black

	if !b.IsStalemate() {
		t.Error("expected stalemate")
	}
	if b.IsCheckmate() {
		t.Error("stalemate position should not be checkmate")
	}
}

func TestCheckmate(t *testing.T) {
	b := NewEmptyBoard()
	b.Set(mustSquare(t, "h8"), Piece{Type: King, Color: Black})
	b.Set(mustSquare(t, "g6"), Piece{Type: King, Color: White})
	b.Set(mustSquare(t, "g7"), Piece{Type: Queen, Color: White})
	b.turn = Black

	if !b.IsCheckmate() {
		t.Error("expected checkmate")
	}
}

func mustSquare(t *testing.T, s string) Square {
	t.Helper()
	sq, err := ParseSquare(s)
	if err != nil {
		t.Fatalf("ParseSquare(%q): %v", s, err)
	}
	return sq
}
