package engine

import "strings"

var backRank = [8]PieceType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}

// Board is a plain 8x8 piece grid with side-to-move tracking. It has no move
// legality logic yet — that lives in the (future) move generator built on
// top of this representation.
type Board struct {
	squares [8][8]Piece
	turn    Color
}

func NewBoard() *Board {
	b := &Board{turn: White}
	for file := 0; file < 8; file++ {
		b.squares[file][1] = Piece{Type: Pawn, Color: White}
		b.squares[file][6] = Piece{Type: Pawn, Color: Black}
		b.squares[file][0] = Piece{Type: backRank[file], Color: White}
		b.squares[file][7] = Piece{Type: backRank[file], Color: Black}
	}
	return b
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

// MakeMove relocates the piece on m.From to m.To without checking legality.
func (b *Board) MakeMove(m Move) {
	p := b.At(m.From)
	if m.Promotion != None {
		p.Type = m.Promotion
	}
	b.Set(m.To, p)
	b.Set(m.From, Piece{})
	b.turn = b.turn.Opponent()
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
