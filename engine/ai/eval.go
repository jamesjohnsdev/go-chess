package ai

import "github.com/jamesjohnsdev/go-chess/engine"

var pieceValues = map[engine.PieceType]int{
	engine.Pawn:   100,
	engine.Knight: 320,
	engine.Bishop: 330,
	engine.Rook:   500,
	engine.Queen:  900,
}

// evaluate scores material from the perspective of the side to move:
// positive favors b.Turn().
func evaluate(b *engine.Board) int {
	score := 0
	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			p := b.At(engine.NewSquare(file, rank))
			if p.Type == engine.None {
				continue
			}
			v := pieceValues[p.Type]
			if p.Color == b.Turn() {
				score += v
			} else {
				score -= v
			}
		}
	}
	return score
}
