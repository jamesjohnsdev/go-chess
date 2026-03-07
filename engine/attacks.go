package engine

var knightOffsets = [8][2]int{{1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}, {-2, -1}, {-2, 1}, {-1, 2}}
var kingOffsets = [8][2]int{{1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}, {0, -1}, {1, -1}}
var bishopDirs = [4][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
var rookDirs = [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

// isAttacked reports whether sq is attacked by any piece of color by.
func (b *Board) isAttacked(sq Square, by Color) bool {
	rankOffset := -1
	if by == Black {
		rankOffset = 1
	}
	for _, fileOffset := range [2]int{-1, 1} {
		from := Square{File: sq.File + fileOffset, Rank: sq.Rank + rankOffset}
		if from.Valid() {
			if p := b.At(from); p.Type == Pawn && p.Color == by {
				return true
			}
		}
	}

	for _, o := range knightOffsets {
		from := Square{File: sq.File + o[0], Rank: sq.Rank + o[1]}
		if from.Valid() {
			if p := b.At(from); p.Type == Knight && p.Color == by {
				return true
			}
		}
	}

	for _, o := range kingOffsets {
		from := Square{File: sq.File + o[0], Rank: sq.Rank + o[1]}
		if from.Valid() {
			if p := b.At(from); p.Type == King && p.Color == by {
				return true
			}
		}
	}

	for _, d := range bishopDirs {
		if b.slideAttacks(sq, d, by, Bishop, Queen) {
			return true
		}
	}
	for _, d := range rookDirs {
		if b.slideAttacks(sq, d, by, Rook, Queen) {
			return true
		}
	}

	return false
}

func (b *Board) slideAttacks(sq Square, dir [2]int, by Color, types ...PieceType) bool {
	cur := Square{File: sq.File + dir[0], Rank: sq.Rank + dir[1]}
	for cur.Valid() {
		p := b.At(cur)
		if p.Type != None {
			if p.Color == by {
				for _, t := range types {
					if p.Type == t {
						return true
					}
				}
			}
			return false
		}
		cur = Square{File: cur.File + dir[0], Rank: cur.Rank + dir[1]}
	}
	return false
}
