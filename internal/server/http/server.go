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
		method  string
		handler http.HandlerFunc
	}{
		{"/items_counts", "GET", httpServer.itemsCountsHandler},
		{"/create", "POST", httpServer.createHandler},
	}

	for _, r := range routes {
		httpServer.serveMux.HandleFunc(r.pattern, methodMiddleware(r.method, r.handler))
	}

	return httpServer
}

func (s *Server) Run(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), s.serveMux)
}

func (s *Server) itemsCountsHandler(w http.ResponseWriter, _ *http.Request) {
	itemsCounts := s.queueManager.ItemsCounts()

	outputData := map[string]any{
		"status":       "done",
		"items_counts": itemsCounts,
		"queues_count": len(itemsCounts),
	}

	writeJson(w, outputData, http.StatusOK)
}

func (s *Server) createHandler(w http.ResponseWriter, r *http.Request) {
	data, err := readJson(r)

	if err != nil {
		writeJson(w, map[string]any{"status": "incorrect json"}, http.StatusBadRequest)

		return
	}

	key, err := getStringValue(data, "key")

	if err != nil {
		writeJson(w, map[string]any{"status": err.Error()}, http.StatusBadRequest)

		return
	}

	isNew := s.queueManager.Create(key)

	if isNew {
		writeJson(w, map[string]any{"status": "the queue has been created"}, http.StatusCreated)

		return
	}

	writeJson(w, map[string]any{"status": "a queue with this key has already been created"}, http.StatusOK)
}

func methodMiddleware(allowedMethod string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			writeJson(w, map[string]any{"status": "method not allowed"}, http.StatusMethodNotAllowed)

			return
		}

		handler(w, r)
	}
}

func writeJson(w http.ResponseWriter, data map[string]any, code int) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
	}
}

func readJson(r *http.Request) (map[string]any, error) {
	var data map[string]any

	return data, json.NewDecoder(r.Body).Decode(&data)
}

func getStringValue(data map[string]any, key string) (string, error) {
	value, ok := data[key]

	if !ok {
		return "", fmt.Errorf("the key \"%v\" is missing", key)
	}

	stringValue, ok := value.(string)

	if !ok || stringValue == "" {
		return "", fmt.Errorf("the value of the key \"%v\" is incorrect", key)
	}

	return stringValue, nil
}
