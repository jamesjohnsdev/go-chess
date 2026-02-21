// Command chess is a terminal chess board printer and scratch REPL. Move
// legality and computer opponents land once the engine's move generator
// exists.
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
		fmt.Printf("%s to move (uci, e.g. e2e4, or 'quit'): ", b.Turn())

		if !scanner.Scan() {
			return
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "quit" {
			return
		}
		if len(input) < 4 {
			fmt.Println("expected a move like e2e4")
			continue
		}

		from, err := engine.ParseSquare(input[0:2])
		if err != nil {
			fmt.Println(err)
			continue
		}
		to, err := engine.ParseSquare(input[2:4])
		if err != nil {
			fmt.Println(err)
			continue
		}
		b.MakeMove(engine.Move{From: from, To: to})
	}
}
