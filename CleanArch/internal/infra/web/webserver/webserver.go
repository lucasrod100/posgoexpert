package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      []MethodHandler
	WebServerPort string
}

type MethodHandler struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      []MethodHandler{},
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(handler MethodHandler) {
	s.Handlers = append(s.Handlers, handler)
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for _, handler := range s.Handlers {
		switch handler.Method {
		case "GET":
			s.Router.Get(handler.Path, handler.Handler)
		case "POST":
			s.Router.Post(handler.Path, handler.Handler)
		case "PUT":
			s.Router.Put(handler.Path, handler.Handler)
		case "DELETE":
			s.Router.Delete(handler.Path, handler.Handler)
		default:
			s.Router.Handle(handler.Path, handler.Handler)
		}
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
