package httpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/anoobz/dualread/auth/internal/store"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type server struct {
	routers *routers
	logger  *log.Logger
	store   store.Store
	port    int
}

type routers struct {
	baseRouter  *mux.Router
	adminRouter *mux.Router
}

func NewServer(store store.Store, logger *log.Logger, port int) *server {
	baseRouter := mux.NewRouter().PathPrefix("/auth").Subrouter()
	adminRouter := baseRouter.PathPrefix("/admin").Subrouter()

	s := &server{
		routers: &routers{
			baseRouter:  baseRouter,
			adminRouter: adminRouter,
		},
		store:  store,
		logger: logger,
		port:   port,
	}

	s.registerRoutes()
	s.configMiddlewares()
	return s
}

func (s *server) Start() error {
	s.logger.Printf("Listening on port: %d", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.routers.baseRouter)
}

func (s *server) registerRoutes() {
	s.routers.baseRouter.HandleFunc("/login", s.login()).Methods("Post")
	s.routers.baseRouter.HandleFunc("/register", s.register()).Methods("Post")
	s.routers.baseRouter.HandleFunc("/refresh-access-token", s.refreshAccessToken()).
		Methods("Post")

	s.routers.adminRouter.HandleFunc("/user/{id:[0-9]+}", s.getUser()).Methods("Get")
	s.routers.adminRouter.HandleFunc("/user-page/{page:[0-9]+}", s.getUserPage()).
		Methods("Get")
	s.routers.adminRouter.HandleFunc("/user", s.getAllUsers()).Methods("Get")
	s.routers.adminRouter.HandleFunc("/user", s.insertUser()).Methods("Post")
	s.routers.adminRouter.HandleFunc("/user/{id:[0-9]+}", s.updateUser()).
		Methods("Post")
	s.routers.adminRouter.HandleFunc("/user/{id:[0-9]+}", s.deleteUser()).
		Methods("Delete")

	s.routers.adminRouter.HandleFunc("/auth-token/{id}", s.getAuthToken()).Methods("Get")
	s.routers.adminRouter.HandleFunc("/auth-token-page/{page:[0-9]+}", s.getAuthTokenPage()).
		Methods("Get")
	s.routers.adminRouter.HandleFunc("/auth-token", s.getAllAuthTokens()).Methods("Get")
	s.routers.adminRouter.HandleFunc("/auth-token/{id}", s.deleteAuthToken()).
		Methods("Delete")
}

func (s *server) configMiddlewares() {
	s.routers.baseRouter.Use(s.logRequest)

	corsOrigin := []string{os.Getenv("CORS_ORIGIN")}
	s.routers.baseRouter.Use(handlers.CORS(handlers.AllowedOrigins(corsOrigin)))

	s.routers.adminRouter.Use(s.validateAccessToken)
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, struct {
		Error string `json:"error"`
	}{Error: err.Error()})
}

func (s *server) respond(
	w http.ResponseWriter, r *http.Request,
	code int, data interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Printf("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		s.logger.Printf(
			"completed with %d: %s, in %v",
			rw.statusCode,
			http.StatusText(rw.statusCode),
			time.Since(start))
	})
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routers.baseRouter.ServeHTTP(w, r)
}
