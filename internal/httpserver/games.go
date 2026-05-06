package httpserver

import (
	"context"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"

	"github.com/jamesjohnsdev/go-chess/engine"
	"github.com/jamesjohnsdev/go-chess/internal/game"
)

type CreateGameInput struct {
	Body *struct {
		Computer string `json:"computer,omitempty" enum:"white,black,none" default:"none" doc:"which side the computer plays; \"none\" (the default) for two human players"`
	}
}

func (in *CreateGameInput) computer() string {
	if in.Body == nil {
		return ""
	}
	return in.Body.Computer
}

type CreateGameOutput struct {
	Body struct {
		ID         string `json:"id"`
		WhiteToken string `json:"white_token,omitempty" doc:"omitted if the computer plays white"`
		BlackToken string `json:"black_token,omitempty" doc:"omitted if the computer plays black"`
		Computer   string `json:"computer,omitempty" doc:"which side the computer plays, if any"`
	}
}

func registerCreateGame(api huma.API, store *game.Store) {
	huma.Register(api, huma.Operation{
		OperationID: "create-game",
		Method:      http.MethodPost,
		Path:        "/games",
		Summary:     "Create a new live game",
	}, func(ctx context.Context, input *CreateGameInput) (*CreateGameOutput, error) {
		out := &CreateGameOutput{}
		switch input.computer() {
		case "white":
			s := store.CreateVsComputer(engine.White)
			out.Body.ID = s.ID
			_, out.Body.BlackToken = s.Tokens()
			out.Body.Computer = "white"
		case "black":
			s := store.CreateVsComputer(engine.Black)
			out.Body.ID = s.ID
			out.Body.WhiteToken, _ = s.Tokens()
			out.Body.Computer = "black"
		default:
			s := store.Create()
			out.Body.ID = s.ID
			out.Body.WhiteToken, out.Body.BlackToken = s.Tokens()
		}
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
// frames: moves, chat, resign, offer_draw, accept_draw, decline_draw, and
// the resulting state updates.
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

			writeErr := func(err error) {
				_ = wsjson.Write(ctx, conn, game.Event{Type: "error", Error: err.Error(), State: session.State()})
			}

			switch msg.Type {
			case "move":
				if msg.Move == nil {
					continue
				}
				mv, err := msg.Move.ToMove()
				if err != nil {
					writeErr(err)
					continue
				}
				if err := session.MakeMove(color, mv); err != nil {
					writeErr(err)
				}
			case "chat":
				session.AddChat(color.String(), msg.Text)
			case "resign":
				if err := session.Resign(color); err != nil {
					writeErr(err)
				}
			case "offer_draw":
				if err := session.OfferDraw(color); err != nil {
					writeErr(err)
				}
			case "accept_draw":
				if err := session.AcceptDraw(color); err != nil {
					writeErr(err)
				}
			case "decline_draw":
				if err := session.DeclineDraw(color); err != nil {
					writeErr(err)
				}
			}
		}
	})
}
