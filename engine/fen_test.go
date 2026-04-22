package engine

import (
	"reflect"
	"strings"
	"testing"
)

func TestFENStartingPosition(t *testing.T) {
	if got := NewBoard().FEN(); got != StartFEN {
		t.Errorf("NewBoard().FEN() = %q, want %q", got, StartFEN)
	}
}

func TestParseFENStartingPosition(t *testing.T) {
	b, err := ParseFEN(StartFEN)
	if err != nil {
		t.Fatalf("ParseFEN(StartFEN): %v", err)
	}
	if !reflect.DeepEqual(*b, *NewBoard()) {
		t.Errorf("ParseFEN(StartFEN) = %+v, want %+v", *b, *NewBoard())
	}
}

func TestFENEnPassantField(t *testing.T) {
	b := NewBoard()
	b.MakeMove(Move{From: mustSquare(t, "e2"), To: mustSquare(t, "e4")})

	fields := strings.Fields(b.FEN())
	if fields[3] != "e3" {
		t.Errorf("en passant field = %q, want e3", fields[3])
	}
}

func TestFENCastlingField(t *testing.T) {
	b := NewEmptyBoard()
	b.Set(mustSquare(t, "e1"), Piece{Type: King, Color: White})
	b.Set(mustSquare(t, "h1"), Piece{Type: Rook, Color: White})
	b.Set(mustSquare(t, "e8"), Piece{Type: King, Color: Black})
	b.castling.WhiteKingside = true

	if fields := strings.Fields(b.FEN()); fields[2] != "K" {
		t.Errorf("castling field = %q, want K", fields[2])
	}

	b.MakeMove(Move{From: mustSquare(t, "e1"), To: mustSquare(t, "g1")})
	if fields := strings.Fields(b.FEN()); fields[2] != "-" {
		t.Errorf("castling field after castling = %q, want -", fields[2])
	}
}

func TestFENRoundTripThroughGame(t *testing.T) {
	b := NewBoard()
	for i := 0; i < 12; i++ {
		moves := b.LegalMoves()
		if len(moves) == 0 {
			break
		}
		b.MakeMove(moves[i%len(moves)])

		fen := b.FEN()
		parsed, err := ParseFEN(fen)
		if err != nil {
			t.Fatalf("ply %d: ParseFEN(%q): %v", i, fen, err)
		}
		if got := parsed.FEN(); got != fen {
			t.Errorf("ply %d: round-trip FEN = %q, want %q", i, got, fen)
		}
		if !reflect.DeepEqual(*parsed, *b) {
			t.Errorf("ply %d: parsed board differs from original:\nparsed:   %+v\noriginal: %+v", i, *parsed, *b)
		}
	}
}

func TestParseFENErrors(t *testing.T) {
	tests := []struct {
		name string
		fen  string
	}{
		{"garbage", "not a fen"},
		{"missing field", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0"},
		{"too few ranks", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP w KQkq - 0 1"},
		{"rank underflow", "rnbqkbnr/pppppppp/7/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
		{"bad digit", "rnbqkbnr/pppppppp/9/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
		{"bad piece char", "rnbqkbnx/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
		{"bad active color", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR x KQkq - 0 1"},
		{"bad castling char", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w XQkq - 0 1"},
		{"bad en passant square", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq z9 0 1"},
		{"negative halfmove", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - -1 1"},
		{"zero fullmove", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ParseFEN(tt.fen); err == nil {
				t.Errorf("ParseFEN(%q) = nil error, want error", tt.fen)
			}
		})
	}
}
