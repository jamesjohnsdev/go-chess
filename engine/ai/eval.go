package ai

import "github.com/jamesjohnsdev/go-chess/engine"

var pieceValues = map[engine.PieceType]int{
	engine.Pawn:   100,
	engine.Knight: 320,
	engine.Bishop: 330,
	engine.Rook:   500,
	engine.Queen:  900,
}

// evaluate scores material and piece placement from the perspective of the
// side to move: positive favors b.Turn().
func evaluate(b *engine.Board) int {
	score := 0
	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			sq := engine.NewSquare(file, rank)
			p := b.At(sq)
			if p.Type == engine.None {
				continue
			}
			v := pieceValues[p.Type] + pstValue(p, sq)
			if p.Color == b.Turn() {
				score += v
			} else {
				score -= v
			}
		}
	}
	return score
}

// pstValue looks up p's piece-square bonus at sq. Tables are stored from
// White's perspective with index 0 = rank 1; Black's rank is mirrored since
// its pieces develop toward rank 1 instead of rank 8.
func pstValue(p engine.Piece, sq engine.Square) int {
	rank := sq.Rank
	if p.Color == engine.Black {
		rank = 7 - rank
	}
	return pieceSquareTables[p.Type][rank][sq.File]
}

var pieceSquareTables = map[engine.PieceType][8][8]int{
	engine.Pawn:   pawnPST,
	engine.Knight: knightPST,
	engine.Bishop: bishopPST,
	engine.Rook:   rookPST,
	engine.Queen:  queenPST,
	engine.King:   kingPST,
}

// The tables below are Tomasz Michniewski's widely used "simplified
// evaluation function" piece-square tables:
// https://www.chessprogramming.org/Simplified_Evaluation_Function
// Row 0 is rank 1 (White's back rank); row 7 is rank 8.

var pawnPST = [8][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{5, 10, 10, -20, -20, 10, 10, 5},
	{5, -5, -10, 0, 0, -10, -5, 5},
	{0, 0, 0, 20, 20, 0, 0, 0},
	{5, 5, 10, 25, 25, 10, 5, 5},
	{10, 10, 20, 30, 30, 20, 10, 10},
	{50, 50, 50, 50, 50, 50, 50, 50},
	{0, 0, 0, 0, 0, 0, 0, 0},
}

var knightPST = [8][8]int{
	{-50, -40, -30, -30, -30, -30, -40, -50},
	{-40, -20, 0, 5, 5, 0, -20, -40},
	{-30, 5, 10, 15, 15, 10, 5, -30},
	{-30, 0, 15, 20, 20, 15, 0, -30},
	{-30, 5, 15, 20, 20, 15, 5, -30},
	{-30, 0, 10, 15, 15, 10, 0, -30},
	{-40, -20, 0, 0, 0, 0, -20, -40},
	{-50, -40, -30, -30, -30, -30, -40, -50},
}

var bishopPST = [8][8]int{
	{-20, -10, -10, -10, -10, -10, -10, -20},
	{-10, 5, 0, 0, 0, 0, 5, -10},
	{-10, 10, 10, 10, 10, 10, 10, -10},
	{-10, 0, 10, 10, 10, 10, 0, -10},
	{-10, 5, 5, 10, 10, 5, 5, -10},
	{-10, 0, 5, 10, 10, 5, 0, -10},
	{-10, 0, 0, 0, 0, 0, 0, -10},
	{-20, -10, -10, -10, -10, -10, -10, -20},
}

var rookPST = [8][8]int{
	{0, 0, 0, 5, 5, 0, 0, 0},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{5, 10, 10, 10, 10, 10, 10, 5},
	{0, 0, 0, 0, 0, 0, 0, 0},
}

var queenPST = [8][8]int{
	{-20, -10, -10, -5, -5, -10, -10, -20},
	{-10, 0, 5, 0, 0, 0, 0, -10},
	{-10, 5, 5, 5, 5, 5, 0, -10},
	{0, 0, 5, 5, 5, 5, 0, -5},
	{-5, 0, 5, 5, 5, 5, 0, -5},
	{-10, 0, 5, 5, 5, 5, 0, -10},
	{-10, 0, 0, 0, 0, 0, 0, -10},
	{-20, -10, -10, -5, -5, -10, -10, -20},
}

// kingPST is middlegame-only: it favors the back rank and castled corners.
// A separate endgame table (which favors an active, centralized king) would
// need game-phase detection this engine doesn't have yet.
var kingPST = [8][8]int{
	{20, 30, 10, 0, 0, 10, 30, 20},
	{20, 20, 0, 0, 0, 0, 20, 20},
	{-10, -20, -20, -20, -20, -20, -20, -10},
	{-20, -30, -30, -40, -40, -30, -30, -20},
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-30, -40, -40, -50, -50, -40, -40, -30},
	{-30, -40, -40, -50, -50, -40, -40, -30},
}
