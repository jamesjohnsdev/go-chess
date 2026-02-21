package engine

type Color uint8

const (
	White Color = iota
	Black
)

func (c Color) Opponent() Color {
	if c == White {
		return Black
	}
	return White
}

func (c Color) String() string {
	if c == White {
		return "white"
	}
	return "black"
}
