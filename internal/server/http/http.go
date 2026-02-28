package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"super-simple-queues/internal/queue"
)

type Http struct {
	queueManager *queue.Manager
	serveMux     *http.ServeMux
}

type route struct {
	pattern string
	method  string
	handler http.HandlerFunc
}

func NewHttp(queueManager *queue.Manager) *Http {
	httpServer := &Http{queueManager, http.NewServeMux()}

	routes := []route{
		{"/messages_counts", "GET", httpServer.messagesCountsHandler},
		{"/create", "POST", httpServer.createHandler},
	}

	for _, r := range routes {
		httpServer.serveMux.HandleFunc(r.pattern, methodMiddleware(r.method, r.handler))
	}

	return httpServer
}

func (h *Http) Run(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%v", port), h.serveMux)
}

func (h *Http) messagesCountsHandler(w http.ResponseWriter, _ *http.Request) {
	messagesCounts := h.queueManager.MessagesCounts()

	outputData := map[string]any{
		"status":          "done",
		"messages_counts": messagesCounts,
		"queues_count":    len(messagesCounts),
	}

	writeJson(w, outputData, 200)
}

func (h *Http) createHandler(w http.ResponseWriter, r *http.Request) {
	data, err := readJson(r)

	if err != nil {
		writeJson(w, map[string]any{"status": "incorrect json"}, 400)

		return
	}

	key, ok := data["key"]

	if !ok {
		writeJson(w, map[string]any{"status": "the key \"key\" is missing"}, 400)

		return
	}

	stringKey, ok := key.(string)

	if !ok || stringKey == "" {
		writeJson(w, map[string]any{"status": "the key \"key\" is incorrect"}, 400)

		return
	}

	h.queueManager.Create(stringKey)

	writeJson(w, map[string]any{"status": "done"}, 201)
}

func methodMiddleware(allowedMethod string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			writeJson(w, map[string]any{"status": "method not allowed"}, 405)

			return
		}

		handler(w, r)
	}
}

func writeJson(w http.ResponseWriter, data map[string]any, code int) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(data)
}

func readJson(r *http.Request) (map[string]any, error) {
	var data map[string]any

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	return data, nil
}
