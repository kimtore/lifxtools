package server

import (
	"github.com/go-chi/chi"
)

func (s *httpserver) Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/{preset}", s.Query)
	r.Put("/{preset}/active", s.Active)
	r.Put("/{preset}/config/{option}", s.Configure)
	return r
}
