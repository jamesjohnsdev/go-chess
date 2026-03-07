package engine

// LegalMoves returns every move the side to move can legally play: pseudo
// legal moves for each piece, filtered to those that don't leave the mover's
// own king in check.
func (b *Board) LegalMoves() []Move {
	pseudo := b.pseudoLegalMoves()
	legal := make([]Move, 0, len(pseudo))
	mover := b.turn
	for _, m := range pseudo {
		clone := *b
		clone.MakeMove(m)
		if !clone.isAttacked(clone.kingSquare(mover), mover.Opponent()) {
			legal = append(legal, m)
		}
	}
	return legal
}

func (b *Board) pseudoLegalMoves() []Move {
	var moves []Move
	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			from := Square{File: file, Rank: rank}
			p := b.At(from)
			if p.Type == None || p.Color != b.turn {
				continue
			}
			switch p.Type {
			case Pawn:
				moves = append(moves, b.pawnMoves(from, p)...)
			case Knight:
				moves = append(moves, b.stepMoves(from, p, knightOffsets[:])...)
			case Bishop:
				moves = append(moves, b.slideMoves(from, p, bishopDirs[:])...)
			case Rook:
				moves = append(moves, b.slideMoves(from, p, rookDirs[:])...)
			case Queen:
				moves = append(moves, b.slideMoves(from, p, bishopDirs[:])...)
				moves = append(moves, b.slideMoves(from, p, rookDirs[:])...)
			case King:
				moves = append(moves, b.stepMoves(from, p, kingOffsets[:])...)
				moves = append(moves, b.castlingMoves(from, p)...)
			}
		}
	}
	return moves
}

func (b *Board) pawnMoves(from Square, p Piece) []Move {
	var moves []Move
	dir, startRank, promoRank := 1, 1, 7
	if p.Color == Black {
		dir, startRank, promoRank = -1, 6, 0
	}

	one := Square{File: from.File, Rank: from.Rank + dir}
	if one.Valid() && b.At(one).Type == None {
		moves = append(moves, expandPromotion(from, one, promoRank)...)
		if from.Rank == startRank {
			two := Square{File: from.File, Rank: from.Rank + 2*dir}
			if b.At(two).Type == None {
				moves = append(moves, Move{From: from, To: two})
			}
		}
	}

	for _, fileOffset := range [2]int{-1, 1} {
		to := Square{File: from.File + fileOffset, Rank: from.Rank + dir}
		if !to.Valid() {
			continue
		}
		if target := b.At(to); target.Type != None {
			if target.Color != p.Color {
				moves = append(moves, expandPromotion(from, to, promoRank)...)
			}
		} else if to == b.epTarget {
			moves = append(moves, Move{From: from, To: to})
		}
	}

	return moves
}

func expandPromotion(from, to Square, promoRank int) []Move {
	if to.Rank != promoRank {
		return []Move{{From: from, To: to}}
	}
	return []Move{
		{From: from, To: to, Promotion: Queen},
		{From: from, To: to, Promotion: Rook},
		{From: from, To: to, Promotion: Bishop},
		{From: from, To: to, Promotion: Knight},
	}
}

func (b *Board) stepMoves(from Square, p Piece, offsets []([2]int)) []Move {
	var moves []Move
	for _, o := range offsets {
		to := Square{File: from.File + o[0], Rank: from.Rank + o[1]}
		if !to.Valid() {
			continue
		}
		if target := b.At(to); target.Type == None || target.Color != p.Color {
			moves = append(moves, Move{From: from, To: to})
		}
	}
	return moves
}

func (b *Board) slideMoves(from Square, p Piece, dirs []([2]int)) []Move {
	var moves []Move
	for _, d := range dirs {
		cur := Square{File: from.File + d[0], Rank: from.Rank + d[1]}
		for cur.Valid() {
			target := b.At(cur)
			if target.Type == None {
				moves = append(moves, Move{From: from, To: cur})
				cur = Square{File: cur.File + d[0], Rank: cur.Rank + d[1]}
				continue
			}
			if target.Color != p.Color {
				moves = append(moves, Move{From: from, To: cur})
			}
			break
		}
	}
	return moves
}

func (b *Board) castlingMoves(from Square, p Piece) []Move {
	rank := 0
	kingside, queenside := b.castling.WhiteKingside, b.castling.WhiteQueenside
	if p.Color == Black {
		rank = 7
		kingside, queenside = b.castling.BlackKingside, b.castling.BlackQueenside
	}
	if from != (Square{File: 4, Rank: rank}) {
		return nil
	}
	opp := p.Color.Opponent()

	var moves []Move
	if kingside &&
		b.At(Square{File: 5, Rank: rank}).Type == None &&
		b.At(Square{File: 6, Rank: rank}).Type == None &&
		!b.isAttacked(Square{File: 4, Rank: rank}, opp) &&
		!b.isAttacked(Square{File: 5, Rank: rank}, opp) &&
		!b.isAttacked(Square{File: 6, Rank: rank}, opp) {
		moves = append(moves, Move{From: from, To: Square{File: 6, Rank: rank}})
	}
	if queenside &&
		b.At(Square{File: 3, Rank: rank}).Type == None &&
		b.At(Square{File: 2, Rank: rank}).Type == None &&
		b.At(Square{File: 1, Rank: rank}).Type == None &&
		!b.isAttacked(Square{File: 4, Rank: rank}, opp) &&
		!b.isAttacked(Square{File: 3, Rank: rank}, opp) &&
		!b.isAttacked(Square{File: 2, Rank: rank}, opp) {
		moves = append(moves, Move{From: from, To: Square{File: 2, Rank: rank}})
	}
	return moves
}
