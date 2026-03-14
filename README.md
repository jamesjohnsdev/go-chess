# go-chess

Chess in Go. Work in progress.

- `engine/` — standalone chess engine package (board representation, legal move generation)
- `engine/ai/` — computer opponent (fixed-depth negamax with alpha-beta pruning)
- `cmd/chess/` — terminal client for local games, human or against the computer (`-computer white|black|both|none`)
- `cmd/server/` — API server for live games with chat (frontend lives in a separate repo)
- `internal/game/` — live game session model (board + chat)
- `internal/httpserver/` — chi/Huma HTTP API layer

## Development

```sh
go build ./...
go test ./...
```
