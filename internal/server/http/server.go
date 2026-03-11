package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"super-simple-queues/internal/queue"
)

type Server struct {
	queueManager *queue.Manager
	serveMux     *http.ServeMux
}

func NewServer(queueManager *queue.Manager) *Server {
	httpServer := &Server{queueManager, http.NewServeMux()}

	routes := []struct {
		pattern string
		handler http.HandlerFunc
	}{
		{"POST /queues/{key}", httpServer.createHandler},
		{"GET /queues/{key}", httpServer.getHandler},
		{"GET /queues", httpServer.listHandler},
	}

	for _, r := range routes {
		httpServer.serveMux.HandleFunc(r.pattern, r.handler)
	}

	return httpServer
}

func (s *Server) Run(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), s.serveMux)
}

func (s *Server) createHandler(w http.ResponseWriter, r *http.Request) {
	createdNow := s.queueManager.Create(r.PathValue("key"))

	if !createdNow {
		writeJson(w, map[string]any{"message": "a queue with this key has already been created"}, http.StatusOK)

		return
	}

	writeJson(w, map[string]any{"message": "the queue has been created"}, http.StatusCreated)
}

func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	q, ok := s.queueManager.Get(r.PathValue("key"))

	if !ok {
		writeJson(w, map[string]any{"message": "a queue with this key does not exist"}, http.StatusNotFound)

		return
	}

	writeJson(w, map[string]any{"items_count": q.Count()}, http.StatusOK)
}

func (s *Server) listHandler(w http.ResponseWriter, _ *http.Request) {
	queues := s.queueManager.GetAll()

	queuesCount := len(queues)

	queuesInfo := make(map[string]any, queuesCount)

	for key, q := range queues {
		queuesInfo[key] = map[string]any{"items_count": q.Count()}
	}

	writeJson(w, map[string]any{"queues_info": queuesInfo, "queues_count": queuesCount}, http.StatusOK)
}

func writeJson(w http.ResponseWriter, data map[string]any, code int) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
	}
}
