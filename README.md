# go-chess

Chess in Go. Work in progress.

- `engine/` — standalone chess engine package (board representation; move legality and a computer opponent are not implemented yet)
- `cmd/chess/` — terminal client for local/single-player games
- `cmd/server/` — webserver for live games with chat
- `internal/game/` — live game session model (board + chat)
- `internal/httpserver/` — HTTP layer for the webserver
- `web/` — web client (placeholder)

## Development

```sh
go build ./...
go test ./...
```
