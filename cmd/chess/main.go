// Command chess is a terminal client for local games, against another
// human or the built-in computer opponent.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jamesjohnsdev/go-chess/engine"
	"github.com/jamesjohnsdev/go-chess/engine/ai"
)

func main() {
	computer := flag.String("computer", "black", `which side the computer plays: "white", "black", "both", or "none"`)
	flag.Parse()

	b := engine.NewBoard()
	bot := ai.New()
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

		if isComputerTurn(b.Turn(), *computer) {
			move, ok := bot.BestMove(b)
			if !ok {
				fmt.Println("computer has no legal move")
				return
			}
			fmt.Printf("%s to move%s: computer plays %s\n", b.Turn(), status, move)
			b.MakeMove(move)
			continue
		}

		fmt.Printf("%s to move%s (uci, e.g. e2e4, or 'quit'): ", b.Turn(), status)
		if !scanner.Scan() {
			return
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "quit" {
			return
		}

		m, err := engine.ParseMove(input)
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

func isComputerTurn(turn engine.Color, computer string) bool {
	switch computer {
	case "white":
		return turn == engine.White
	case "black":
		return turn == engine.Black
	case "both":
		return true
	default:
		return false
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
