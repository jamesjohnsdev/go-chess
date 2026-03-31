package engine

import "strings"

func (b *Board) HalfmoveClock() int {
	return b.halfmove
}

// IsFiftyMoveDraw reports whether fifty full moves (100 halfmoves) have
// passed without a pawn move or capture.
func (b *Board) IsFiftyMoveDraw() bool {
	return b.halfmove >= 100
}

// IsInsufficientMaterial reports whether neither side has enough material
// to checkmate: bare kings, king+minor vs king, or same-colored-bishop vs
// same-colored-bishop.
func (b *Board) IsInsufficientMaterial() bool {
	type minor struct {
		Type  PieceType
		Color Color
		Sq    Square
	}
	var minors []minor
	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			p := b.squares[file][rank]
			switch p.Type {
			case Pawn, Rook, Queen:
				return false
			case Bishop, Knight:
				minors = append(minors, minor{p.Type, p.Color, Square{File: file, Rank: rank}})
			}
		}
	}

	switch len(minors) {
	case 0, 1:
		return true
	case 2:
		if minors[0].Color == minors[1].Color || minors[0].Type != Bishop || minors[1].Type != Bishop {
			return false
		}
		return squareColor(minors[0].Sq) == squareColor(minors[1].Sq)
	default:
		return false
	}
}

func squareColor(sq Square) int {
	return (sq.File + sq.Rank) % 2
}

// PositionKey returns a value that's equal for two boards iff they have the
// same piece placement, side to move, castling rights, and en passant
// target — everything FIDE's threefold-repetition rule considers part of
// "the same position". Callers track repetition by counting keys over the
// course of a game; Board itself holds no history.
func (b *Board) PositionKey() string {
	var sb strings.Builder
	sb.Grow(64 + 5)
	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			p := b.squares[file][rank]
			sb.WriteByte(byte(p.Type)<<1 | byte(p.Color))
		}
	}
	sb.WriteByte(byte(b.turn))
	sb.WriteByte(boolByte(b.castling.WhiteKingside))
	sb.WriteByte(boolByte(b.castling.WhiteQueenside))
	sb.WriteByte(boolByte(b.castling.BlackKingside))
	sb.WriteByte(boolByte(b.castling.BlackQueenside))
	sb.WriteByte(byte(b.epTarget.File + 1))
	return sb.String()
}

func boolByte(v bool) byte {
	if v {
		return 1
	}
	return 0
}
