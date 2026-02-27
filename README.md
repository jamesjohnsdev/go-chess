# go-chess

Chess in Go. Work in progress.

- `engine/` — standalone chess engine package (board representation; move legality and a computer opponent are not implemented yet)
- `cmd/chess/` — terminal client for local/single-player games
- `cmd/server/` — API server for live games with chat (frontend lives in a separate repo)
- `internal/game/` — live game session model (board + chat)
- `internal/httpserver/` — chi/Huma HTTP API layer

## Development

```sh
go build ./...
go test ./...
```
