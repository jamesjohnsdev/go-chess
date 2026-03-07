// Command chess is a terminal client for local games. Single-player against
// a computer opponent lands once the engine grows one.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jamesjohnsdev/go-chess/engine"
)

func main() {
	b := engine.NewBoard()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(b)

		if b.IsCheckmate() {
			fmt.Printf("checkmate — %s wins\n", b.Turn().Opponent())
			return
		}
		if b.IsStalemate() {
			fmt.Println("stalemate")
			return
		}
		status := ""
		if b.InCheck() {
			status = " (check)"
		}
		fmt.Printf("%s to move%s (uci, e.g. e2e4, or 'quit'): ", b.Turn(), status)

		if !scanner.Scan() {
			return
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "quit" {
			return
		}

		m, err := parseMove(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if !isLegal(b, m) {
			fmt.Println("illegal move")
			continue
		}
		b.MakeMove(m)
	}
}

func parseMove(input string) (engine.Move, error) {
	if len(input) < 4 {
		return engine.Move{}, fmt.Errorf("expected a move like e2e4")
	}
	from, err := engine.ParseSquare(input[0:2])
	if err != nil {
		return engine.Move{}, err
	}
	to, err := engine.ParseSquare(input[2:4])
	if err != nil {
		return engine.Move{}, err
	}
	m := engine.Move{From: from, To: to}
	if len(input) >= 5 {
		m.Promotion = promotionFromLetter(input[4])
	}
	return m, nil
}

func promotionFromLetter(c byte) engine.PieceType {
	switch c {
	case 'q':
		return engine.Queen
	case 'r':
		return engine.Rook
	case 'b':
		return engine.Bishop
	case 'n':
		return engine.Knight
	default:
		return engine.None
	}
}

func isLegal(b *engine.Board, m engine.Move) bool {
	for _, legal := range b.LegalMoves() {
		if legal == m {
			return true
		}
	}
	return false
}
