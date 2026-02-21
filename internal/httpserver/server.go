// Package httpserver wires up the HTTP handlers for the web version of the
// game: static assets now, live-game and chat endpoints (likely WebSocket)
// once internal/game grows session management.
package httpserver

import "net/http"

func New() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealthz)
	mux.Handle("/", http.FileServer(http.Dir("web/static")))
	return mux
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
