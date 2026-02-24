package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"super-simple-queues/internal/queue"
)

type Web struct {
	manager  *queue.Manager
	serveMux *http.ServeMux
}

func NewWeb(manager *queue.Manager) *Web {
	web := &Web{manager, http.NewServeMux()}

	web.serveMux.HandleFunc("/list", web.listHandler)
	web.serveMux.HandleFunc("/create", web.createHandler)

	return web
}

func (w *Web) Run(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%v", port), w.serveMux)
}

func (w *Web) listHandler(writer http.ResponseWriter, req *http.Request) {
	everything := w.manager.GetEverything()

	response := map[string]interface{}{
		"queues":   everything,
		"quantity": len(everything),
	}

	b, err := json.Marshal(response)

	if err != nil {
		log.Println(err)
	}

	_, err = writer.Write(b)

	if err != nil {
		log.Println(err)
	}
}

func (w *Web) createHandler(writer http.ResponseWriter, req *http.Request) {
	qk := req.URL.Query()["key"]

	w.manager.CreateQueue(qk[0])
}
