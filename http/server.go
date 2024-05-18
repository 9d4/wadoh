package http

import (
	"context"
	"io/fs"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	"github.com/9d4/wadoh/html"
	"github.com/9d4/wadoh/storage"
)

type Config struct {
	Address   string `koanf:"address"`
	JWTSecret []byte `koanf:"jwt_secret"`
}

type Server struct {
	config *Config
	ln     net.Listener
	router chi.Router
	server *http.Server

	storage   *storage.Storage
	staticFs  fs.FS
	templates *html.Templates

	tokenAuth *jwtauth.JWTAuth
}

type Option func(*Config)

func NewServer(storage *storage.Storage, opts ...Option) *Server {
	s := &Server{
		config:    &Config{},
		server:    &http.Server{},
		router:    chi.NewRouter(),
		storage:   storage,
		staticFs:  html.StaticFs(),
		templates: html.NewTemplates(),
	}

	for _, fn := range opts {
		if fn != nil {
			fn(s.config)
		}
	}

	s.tokenAuth = jwtauth.New("HS256", s.config.JWTSecret, nil)

	s.router.Use(jwtauth.Verifier(s.tokenAuth))
	initializeRoutes(s)

	s.server.Handler = http.HandlerFunc(s.serveHTTP)
	return s
}

func (s *Server) Serve() (err error) {
	s.ln, err = net.Listen("tcp", s.config.Address)
	if err != nil {
		return err
	}

	go s.server.Serve(s.ln)
	return nil
}

func (s *Server) ShutDown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Server) Address() string {
	return s.ln.Addr().String()
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.Get("/static/*", staticHandler("/static/", s.staticFs))
	s.router.ServeHTTP(w, r)
}

func staticHandler(prefix string, fs fs.FS) func(w http.ResponseWriter, r *http.Request) {
	h := http.StripPrefix(prefix, http.FileServer(http.FS(fs)))
	return h.ServeHTTP
}
