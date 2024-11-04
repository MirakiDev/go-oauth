package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(
		cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowCredentials: true,
		},
	))
	r.Use(middleware.Logger)

	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)
	r.Get("/auth/{provider}", s.beginAuthProviderCallback)
	r.Get("/auth/logout/{provider}", s.logout)
	r.Get("/auth/{provider}/callback", s.getAuthCallbackFunction)
	r.Get("/auth/me", s.getMeHandler)
	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) getAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, "Error in authentication:", err)
		return
	}

	session, _ := gothic.Store.Get(r, "gothic-session")
	session.Values["user"] = user
	session.Save(r, w)

	http.Redirect(w, r, "http://localhost:5173?success=ok", http.StatusFound)
}

func (s *Server) beginAuthProviderCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	fmt.Println("Provider:", provider)

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	gothic.BeginAuthHandler(w, r)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)

	session, err := gothic.Store.Get(r, "gothic-session")

	if err == nil {
		delete(session.Values, "user")
		session.Save(r, w)
	}

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) getMeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := gothic.Store.Get(r, "gothic-session")
	if err != nil {
		http.Error(w, "Error retrieving session", http.StatusUnauthorized)
		return
	}

	user, ok := session.Values["user"].(goth.User)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	jsonResp, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}
