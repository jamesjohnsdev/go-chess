package engine

import "strings"

var backRank = [8]PieceType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}

// noEnPassant marks the absence of an en passant target square. It is never
// a real target since rank 0 can't hold one.
var noEnPassant = Square{File: -1, Rank: -1}

type CastlingRights struct {
	WhiteKingside  bool
	WhiteQueenside bool
	BlackKingside  bool
	BlackQueenside bool
}

// Board is an 8x8 piece grid plus the state needed for legal move
// generation and draw detection: side to move, castling rights, the en
// passant target square, and the fifty-move-rule halfmove clock.
type Board struct {
	squares  [8][8]Piece
	turn     Color
	castling CastlingRights
	epTarget Square
	halfmove int
	fullmove int
}

func NewBoard() *Board {
	b := NewEmptyBoard()
	for file := 0; file < 8; file++ {
		b.squares[file][1] = Piece{Type: Pawn, Color: White}
		b.squares[file][6] = Piece{Type: Pawn, Color: Black}
		b.squares[file][0] = Piece{Type: backRank[file], Color: White}
		b.squares[file][7] = Piece{Type: backRank[file], Color: Black}
	}
	b.castling = CastlingRights{true, true, true, true}
	return b
}

// NewEmptyBoard returns a board with no pieces, White to move, no castling
// rights, and no en passant target. Useful for constructing test positions.
func NewEmptyBoard() *Board {
	return &Board{turn: White, epTarget: noEnPassant, fullmove: 1}
}

func (b *Board) At(sq Square) Piece {
	return b.squares[sq.File][sq.Rank]
}

func (b *Board) Set(sq Square, p Piece) {
	b.squares[sq.File][sq.Rank] = p
}

func (b *Board) Turn() Color {
	return b.turn
}

// FullmoveNumber starts at 1 and increments after each Black move, per FEN
// convention.
func (b *Board) FullmoveNumber() int {
	return b.fullmove
}

func (b *Board) Castling() CastlingRights {
	return b.castling
}

// EnPassant returns the current en passant target square and whether one is
// set.
func (b *Board) EnPassant() (Square, bool) {
	if b.epTarget == noEnPassant {
		return Square{}, false
	}
	return b.epTarget, true
}

// MakeMove applies m without checking legality; callers should only pass
// moves from LegalMoves. It handles captures, en passant, castling
// (including rook relocation), promotion, and castling-rights/en-passant
// bookkeeping.
func (b *Board) MakeMove(m Move) {
	p := b.At(m.From)
	mover := p
	isCapture := b.At(m.To).Type != None || (p.Type == Pawn && m.To == b.epTarget)

	if p.Type == Pawn && m.To == b.epTarget {
		b.Set(Square{File: m.To.File, Rank: m.From.Rank}, Piece{})
	}

	if p.Type == King && abs(m.To.File-m.From.File) == 2 {
		rank := m.From.Rank
		if m.To.File == 6 {
			b.Set(Square{File: 5, Rank: rank}, b.At(Square{File: 7, Rank: rank}))
			b.Set(Square{File: 7, Rank: rank}, Piece{})
		} else {
			b.Set(Square{File: 3, Rank: rank}, b.At(Square{File: 0, Rank: rank}))
			b.Set(Square{File: 0, Rank: rank}, Piece{})
		}
	}

	if m.Promotion != None {
		p.Type = m.Promotion
	}
	b.Set(m.To, p)
	b.Set(m.From, Piece{})

	b.updateCastlingRights(m, mover)

	if p.Type == Pawn && abs(m.To.Rank-m.From.Rank) == 2 {
		b.epTarget = Square{File: m.From.File, Rank: (m.From.Rank + m.To.Rank) / 2}
	} else {
		b.epTarget = noEnPassant
	}

	if mover.Type == Pawn || isCapture {
		b.halfmove = 0
	} else {
		b.halfmove++
	}

	if b.turn == Black {
		b.fullmove++
	}
	b.turn = b.turn.Opponent()
}

func (b *Board) updateCastlingRights(m Move, mover Piece) {
	if mover.Type == King {
		if mover.Color == White {
			b.castling.WhiteKingside = false
			b.castling.WhiteQueenside = false
		} else {
			b.castling.BlackKingside = false
			b.castling.BlackQueenside = false
		}
	}

	clearRookRight := func(sq Square) {
		switch sq {
		case Square{File: 0, Rank: 0}:
			b.castling.WhiteQueenside = false
		case Square{File: 7, Rank: 0}:
			b.castling.WhiteKingside = false
		case Square{File: 0, Rank: 7}:
			b.castling.BlackQueenside = false
		case Square{File: 7, Rank: 7}:
			b.castling.BlackKingside = false
		}
	}
	clearRookRight(m.From)
	clearRookRight(m.To)
}

func (b *Board) kingSquare(c Color) Square {
	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			p := b.squares[file][rank]
			if p.Type == King && p.Color == c {
				return Square{File: file, Rank: rank}
			}
		}
	}
	return Square{File: -1, Rank: -1}
}

func (b *Board) InCheck() bool {
	return b.isAttacked(b.kingSquare(b.turn), b.turn.Opponent())
}

func (b *Board) IsCheckmate() bool {
	return b.InCheck() && len(b.LegalMoves()) == 0
}

func (b *Board) IsStalemate() bool {
	return !b.InCheck() && len(b.LegalMoves()) == 0
}

func (b *Board) String() string {
	var sb strings.Builder
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			sb.WriteString(b.squares[file][rank].String())
			sb.WriteByte(' ')
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
