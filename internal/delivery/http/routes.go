package deliveryHttp

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/pedroctbd/WhatsappClone/internal/chat"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Application struct {
	Logger      *log.Logger
	Hub         *chat.Hub
	ChatService *chat.ChatService
	ServerID    string
}

func (app *Application) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", serveHome)
	r.Get("/ws/{userId}", app.handleConnections)
	return r
}
