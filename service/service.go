package service

import (
	"encoding/json"
	"net/http"

	"github.com/funkytennisball/hera/common"
)

// Service provides an HTTP service
type Service struct {
	Protocol common.Protocol
}

func (s *Service) postHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: error handling
	w.Write([]byte(s.Protocol.Get(common.Key{K: m["key"]}).V))
}

func (s *Service) putHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: error handling
	if err := s.Protocol.Insert(common.Key{K: m["key"]}, common.Entity{V: m["value"]}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) deleteHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: error handling
	if err := s.Protocol.Remove(common.Key{K: m["key"]}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) defaultHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.postHandler(w, r)
	case http.MethodPut:
		s.putHandler(w, r)
	case http.MethodDelete:
		s.deleteHandler(w, r)
	}
}

// Start the instance and the HTTP web server
func (s *Service) Start() {
	http.HandleFunc("/", s.defaultHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
