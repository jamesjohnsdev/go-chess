package engine

import "fmt"

type Move struct {
	From      Square
	To        Square
	Promotion PieceType
}

func (m Move) String() string {
	s := m.From.String() + m.To.String()
	s += m.Promotion.Letter()
	return s
}

// ParseMove parses UCI-style move notation, e.g. "e2e4" or "e7e8q".
func ParseMove(s string) (Move, error) {
	if len(s) < 4 {
		return Move{}, fmt.Errorf("invalid move %q", s)
	}
	from, err := ParseSquare(s[0:2])
	if err != nil {
		return Move{}, err
	}
	to, err := ParseSquare(s[2:4])
	if err != nil {
		return Move{}, err
	}
	m := Move{From: from, To: to}
	if len(s) >= 5 {
		m.Promotion = promotionFromLetter(s[4])
		if m.Promotion == None {
			return Move{}, fmt.Errorf("invalid promotion %q", s[4:5])
		}
	}
	return m, nil
}

func promotionFromLetter(c byte) PieceType {
	switch c {
	case 'q':
		return Queen
	case 'r':
		return Rook
	case 'b':
		return Bishop
	case 'n':
		return Knight
	default:
		return None
	}
}
