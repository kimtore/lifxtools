package server

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/dorkowscy/lifxtool/pkg/runner"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

// Implements a HTTP server that controls the effect manager

type httpserver struct {
	manager runner.Manager
}

// Data returned on GET /{preset}
type effectState struct {
	Description string                 `json:"description"`
	Active      bool                   `json:"active"`
	Config      map[string]interface{} `json:"config"`
}

func New(manager runner.Manager) *httpserver {
	return &httpserver{
		manager: manager,
	}
}

func (s *httpserver) Query(w http.ResponseWriter, r *http.Request) {
	preset := chi.URLParam(r, "preset")

	if !s.manager.Exists(preset) {
		http.NotFound(w, r)
		return
	}

	active := s.manager.IsActive(preset)
	cfg := s.manager.Configuration(preset)

	log.Debugf("Query(%s): active=%v, cfg=%v", preset, active, cfg)

	w.Header().Set("content-type", "application/json")

	json.NewEncoder(w).Encode(&effectState{
		Active: active,
		Config: cfg,
	})
}

func (s *httpserver) Active(w http.ResponseWriter, r *http.Request) {
	preset := chi.URLParam(r, "preset")

	defer r.Body.Close()
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	active, err := strconv.ParseBool(string(payload))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if active {
		err = s.manager.StartPreset(preset)
	} else {
		err = s.manager.StopPreset(preset)
	}
	if err == nil {
		return
	}
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func (s *httpserver) Configure(w http.ResponseWriter, r *http.Request) {
	preset := chi.URLParam(r, "preset")
	option := chi.URLParam(r, "option")

	defer r.Body.Close()

	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		log.Errorf("Set %s.%s: %s", preset, option, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	value := unmarshal(buf.String())

	log.Debugf("Set %T %s.%s=%v", value, preset, option, value)

	cfg := map[string]interface{}{
		option: value,
	}

	err = s.manager.Configure(preset, cfg)
	if err != nil {
		log.Errorf("Configuration failed: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	log.Debugf("cfg=%v", s.manager.Configuration(preset))
}

func unmarshal(value string) interface{} {
	f, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return f
	}

	b, err := strconv.ParseBool(value)
	if err == nil {
		return b
	}

	d, err := time.ParseDuration(value)
	if err == nil {
		return d
	}

	return value
}
