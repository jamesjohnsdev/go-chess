package engine

import "fmt"

// Square uses 0-based File (a-h) and Rank (1-8) coordinates.
type Square struct {
	File int
	Rank int
}

func NewSquare(file, rank int) Square {
	return Square{File: file, Rank: rank}
}

func (s Square) Valid() bool {
	return s.File >= 0 && s.File < 8 && s.Rank >= 0 && s.Rank < 8
}

func (s Square) String() string {
	if !s.Valid() {
		return "-"
	}
	return fmt.Sprintf("%c%d", 'a'+s.File, s.Rank+1)
}

func ParseSquare(str string) (Square, error) {
	if len(str) != 2 {
		return Square{}, fmt.Errorf("invalid square %q", str)
	}
	file := int(str[0] - 'a')
	rank := int(str[1] - '1')
	sq := Square{File: file, Rank: rank}
	if !sq.Valid() {
		return Square{}, fmt.Errorf("invalid square %q", str)
	}
	return sq, nil
}
