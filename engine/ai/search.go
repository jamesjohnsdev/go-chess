// Package ai implements a computer opponent on top of the engine package's
// legal move generator: fixed-depth negamax with alpha-beta pruning and a
// material-only evaluation.
package ai

import "github.com/jamesjohnsdev/go-chess/engine"

const (
	// inf and mateScore stay well clear of int overflow on negation, unlike
	// math.MinInt/MaxInt.
	inf       = 1 << 30
	mateScore = 1 << 20
)

const DefaultDepth = 3

type Engine struct {
	// Depth is the number of half-moves (plies) to search. Values <= 0 fall
	// back to DefaultDepth.
	Depth int
}

func New() *Engine {
	return &Engine{Depth: DefaultDepth}
}

func (e *Engine) depth() int {
	if e.Depth <= 0 {
		return DefaultDepth
	}
	return e.Depth
}

// BestMove returns the best legal move for the side to move and true, or
// the zero Move and false if the game has already ended.
func (e *Engine) BestMove(b *engine.Board) (engine.Move, bool) {
	moves := b.LegalMoves()
	if len(moves) == 0 {
		return engine.Move{}, false
	}

	depth := e.depth()
	best := moves[0]
	bestScore := -inf
	alpha, beta := -inf, inf
	for _, m := range moves {
		clone := *b
		clone.MakeMove(m)
		score := -negamax(&clone, depth-1, -beta, -alpha)
		if score > bestScore {
			bestScore = score
			best = m
		}
		if score > alpha {
			alpha = score
		}
	}
	return best, true
}

func negamax(b *engine.Board, depth, alpha, beta int) int {
	moves := b.LegalMoves()
	if len(moves) == 0 {
		if b.InCheck() {
			// Favor faster mates: a mate found with more depth left
			// happened sooner in the game tree.
			return -(mateScore + depth)
		}
		return 0
	}
	if depth == 0 {
		return evaluate(b)
	}

	best := -inf
	for _, m := range moves {
		clone := *b
		clone.MakeMove(m)
		score := -negamax(&clone, depth-1, -beta, -alpha)
		if score > best {
			best = score
		}
		if score > alpha {
			alpha = score
		}
		if alpha >= beta {
			break
		}
	}
	return best
}
