package httpserver

import (
	"context"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"

	"github.com/jamesjohnsdev/go-chess/internal/game"
)

type CreateGameOutput struct {
	Body struct {
		ID         string `json:"id"`
		WhiteToken string `json:"white_token"`
		BlackToken string `json:"black_token"`
	}
}

func registerCreateGame(api huma.API, store *game.Store) {
	huma.Register(api, huma.Operation{
		OperationID: "create-game",
		Method:      http.MethodPost,
		Path:        "/games",
		Summary:     "Create a new live game",
	}, func(ctx context.Context, input *struct{}) (*CreateGameOutput, error) {
		s := store.Create()
		out := &CreateGameOutput{}
		out.Body.ID = s.ID
		out.Body.WhiteToken, out.Body.BlackToken = s.Tokens()
		return out, nil
	})
}

type GetGameInput struct {
	ID string `path:"id"`
}

type GetGameOutput struct {
	Body game.State
}

func registerGetGame(api huma.API, store *game.Store) {
	huma.Register(api, huma.Operation{
		OperationID: "get-game",
		Method:      http.MethodGet,
		Path:        "/games/{id}",
		Summary:     "Get live game state",
	}, func(ctx context.Context, input *GetGameInput) (*GetGameOutput, error) {
		s, ok := store.Get(input.ID)
		if !ok {
			return nil, huma.Error404NotFound("game not found")
		}
		return &GetGameOutput{Body: s.State()}, nil
	})
}

// registerGameSocket wires /games/{id}/ws directly on the chi router: Huma
// doesn't manage hijacked connections, so this bypasses it. A client
// authenticates as a color with ?token=<white_token|black_token> from
// CreateGameOutput, then exchanges game.ClientMessage/game.Event JSON
// frames for moves, chat, and state updates.
func registerGameSocket(router chi.Router, store *game.Store) {
	router.Get("/games/{id}/ws", func(w http.ResponseWriter, r *http.Request) {
		session, ok := store.Get(chi.URLParam(r, "id"))
		if !ok {
			http.Error(w, "game not found", http.StatusNotFound)
			return
		}
		color, err := session.ColorForToken(r.URL.Query().Get("token"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		// InsecureSkipVerify/OriginPatterns will be needed once the
		// frontend is served from a different origin than this API.
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()

		ctx := r.Context()
		events, cancel := session.Subscribe()
		defer cancel()

		go func() {
			for e := range events {
				if wsjson.Write(ctx, conn, e) != nil {
					return
				}
			}
		}()

		for {
			var msg game.ClientMessage
			if err := wsjson.Read(ctx, conn, &msg); err != nil {
				return
			}

			switch msg.Type {
			case "move":
				if msg.Move == nil {
					continue
				}
				mv, err := msg.Move.ToMove()
				if err != nil {
					_ = wsjson.Write(ctx, conn, game.Event{Type: "error", Error: err.Error(), State: session.State()})
					continue
				}
				if err := session.MakeMove(color, mv); err != nil {
					_ = wsjson.Write(ctx, conn, game.Event{Type: "error", Error: err.Error(), State: session.State()})
				}
			case "chat":
				session.AddChat(color.String(), msg.Text)
			}
		}
	})
}
