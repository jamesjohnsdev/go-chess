package engine

type Move struct {
	From      Square
	To        Square
	Promotion PieceType
}

func (m Move) String() string {
	s := m.From.String() + m.To.String()
	if letter, ok := pieceLetters[m.Promotion]; ok {
		s += letter
	}
	return s
}
