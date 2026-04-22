# go-chess

Chess in Go. Work in progress.

- `engine/` — standalone chess engine package (board representation, legal move generation, checkmate/stalemate/fifty-move-rule/insufficient-material/threefold-repetition detection, FEN import/export)
- `engine/ai/` — computer opponent (fixed-depth negamax with alpha-beta pruning)
- `cmd/chess/` — terminal client for local games, human or against the computer (`-computer white|black|both|none`, `-fen "..."` to start from a position, `fen` at the prompt to print the current one)
- `cmd/server/` — API server for live games with chat (frontend lives in a separate repo)
- `internal/game/` — live game session store: board, chat, per-session pub-sub event stream, optional built-in computer opponent for one side
- `internal/httpserver/` — chi/Huma REST API (`POST /games` with an optional `{"computer": "white"|"black"}` body, `GET /games/{id}` — state includes both an ASCII board and FEN) plus a `/games/{id}/ws` WebSocket for moves and chat

## Development

```sh
go build ./...
go test ./...
```
