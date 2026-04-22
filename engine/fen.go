package engine

import (
	"fmt"
	"strconv"
	"strings"
)

// StartFEN is the FEN for the standard starting position.
const StartFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// FEN returns the Forsyth-Edwards Notation for the current position.
func (b *Board) FEN() string {
	var sb strings.Builder

	for rank := 7; rank >= 0; rank-- {
		empty := 0
		for file := 0; file < 8; file++ {
			p := b.squares[file][rank]
			if p.Type == None {
				empty++
				continue
			}
			if empty > 0 {
				sb.WriteString(strconv.Itoa(empty))
				empty = 0
			}
			sb.WriteString(p.String())
		}
		if empty > 0 {
			sb.WriteString(strconv.Itoa(empty))
		}
		if rank > 0 {
			sb.WriteByte('/')
		}
	}

	sb.WriteByte(' ')
	if b.turn == White {
		sb.WriteByte('w')
	} else {
		sb.WriteByte('b')
	}

	sb.WriteByte(' ')
	sb.WriteString(castlingFEN(b.castling))

	sb.WriteByte(' ')
	if ep, ok := b.EnPassant(); ok {
		sb.WriteString(ep.String())
	} else {
		sb.WriteByte('-')
	}

	fmt.Fprintf(&sb, " %d %d", b.halfmove, b.fullmove)
	return sb.String()
}

func castlingFEN(c CastlingRights) string {
	s := ""
	if c.WhiteKingside {
		s += "K"
	}
	if c.WhiteQueenside {
		s += "Q"
	}
	if c.BlackKingside {
		s += "k"
	}
	if c.BlackQueenside {
		s += "q"
	}
	if s == "" {
		return "-"
	}
	return s
}

// ParseFEN parses Forsyth-Edwards Notation into a Board.
func ParseFEN(fen string) (*Board, error) {
	fields := strings.Fields(fen)
	if len(fields) != 6 {
		return nil, fmt.Errorf("invalid FEN %q: expected 6 fields, got %d", fen, len(fields))
	}

	b := &Board{epTarget: noEnPassant}

	if err := parsePlacementFEN(b, fields[0]); err != nil {
		return nil, fmt.Errorf("invalid FEN %q: %w", fen, err)
	}

	switch fields[1] {
	case "w":
		b.turn = White
	case "b":
		b.turn = Black
	default:
		return nil, fmt.Errorf("invalid FEN %q: bad active color %q", fen, fields[1])
	}

	if fields[2] != "-" {
		for _, c := range fields[2] {
			switch c {
			case 'K':
				b.castling.WhiteKingside = true
			case 'Q':
				b.castling.WhiteQueenside = true
			case 'k':
				b.castling.BlackKingside = true
			case 'q':
				b.castling.BlackQueenside = true
			default:
				return nil, fmt.Errorf("invalid FEN %q: bad castling rights %q", fen, fields[2])
			}
		}
	}

	if fields[3] != "-" {
		sq, err := ParseSquare(fields[3])
		if err != nil {
			return nil, fmt.Errorf("invalid FEN %q: bad en passant target: %w", fen, err)
		}
		b.epTarget = sq
	}

	halfmove, err := strconv.Atoi(fields[4])
	if err != nil || halfmove < 0 {
		return nil, fmt.Errorf("invalid FEN %q: bad halfmove clock %q", fen, fields[4])
	}
	b.halfmove = halfmove

	fullmove, err := strconv.Atoi(fields[5])
	if err != nil || fullmove < 1 {
		return nil, fmt.Errorf("invalid FEN %q: bad fullmove number %q", fen, fields[5])
	}
	b.fullmove = fullmove

	return b, nil
}

func parsePlacementFEN(b *Board, placement string) error {
	ranks := strings.Split(placement, "/")
	if len(ranks) != 8 {
		return fmt.Errorf("expected 8 ranks, got %d", len(ranks))
	}
	for i, rankStr := range ranks {
		rank := 7 - i
		file := 0
		for _, c := range rankStr {
			if c >= '1' && c <= '8' {
				file += int(c - '0')
				continue
			}
			if file >= 8 {
				return fmt.Errorf("rank %d overflows", rank+1)
			}
			p, err := pieceFromFEN(c)
			if err != nil {
				return err
			}
			b.squares[file][rank] = p
			file++
		}
		if file != 8 {
			return fmt.Errorf("rank %d has %d squares, want 8", rank+1, file)
		}
	}
	return nil
}

func pieceFromFEN(c rune) (Piece, error) {
	color := White
	lower := c
	if c >= 'a' && c <= 'z' {
		color = Black
	} else {
		lower = c + ('a' - 'A')
	}
	switch lower {
	case 'p':
		return Piece{Type: Pawn, Color: color}, nil
	case 'n':
		return Piece{Type: Knight, Color: color}, nil
	case 'b':
		return Piece{Type: Bishop, Color: color}, nil
	case 'r':
		return Piece{Type: Rook, Color: color}, nil
	case 'q':
		return Piece{Type: Queen, Color: color}, nil
	case 'k':
		return Piece{Type: King, Color: color}, nil
	default:
		return Piece{}, fmt.Errorf("bad piece char %q", c)
	}
}
