// Command server runs the web version of go-chess: live games, chat, and
// the static frontend.
package main

import (
	"log"
	"net/http"

	"github.com/jamesjohnsdev/go-chess/internal/httpserver"
)

func main() {
	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, httpserver.New()))
}
