package engine

import "strings"

type PieceType uint8

const (
	None PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

type Piece struct {
	Type  PieceType
	Color Color
}

var pieceLetters = map[PieceType]string{
	Pawn:   "p",
	Knight: "n",
	Bishop: "b",
	Rook:   "r",
	Queen:  "q",
	King:   "k",
}

// Letter returns the lowercase algebraic letter for t, or "" for None.
func (t PieceType) Letter() string {
	return pieceLetters[t]
}

func (p Piece) String() string {
	letter, ok := pieceLetters[p.Type]
	if !ok {
		return "."
	}
	if p.Color == White {
		return strings.ToUpper(letter)
	}
	return letter
}
