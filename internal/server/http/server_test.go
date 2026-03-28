package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"super-simple-queues/internal/queue"
	"testing"
)

var baseURL string

func TestMain(m *testing.M) {
	const queueChunkSize = 1024

	httpServer := NewServer(queue.NewManager(queueChunkSize))

	ts := httptest.NewServer(httpServer.serveMux)

	baseURL = ts.URL

	code := m.Run()

	ts.Close()

	os.Exit(code)
}

func TestServer_createHandler(t *testing.T) {
	statusCode, response := request(t, "POST", "/queues/test")
	t.Cleanup(func() { _, _ = request(t, "DELETE", "/queues/test") })

	checkStatusCode(t, statusCode, http.StatusCreated)
	checkResponse(t, response, map[string]any{"message": "the queue has been created"})

	statusCode, response = request(t, "POST", "/queues/test")

	checkStatusCode(t, statusCode, http.StatusOK)
	checkResponse(t, response, map[string]any{"message": "a queue with this key has already been created"})
}

func TestServer_getHandler(t *testing.T) {
	statusCode, response := request(t, "GET", "/queues/test")

	checkStatusCode(t, statusCode, http.StatusNotFound)
	checkResponse(t, response, map[string]any{"message": "a queue with this key does not exist"})

	_, _ = request(t, "POST", "/queues/test")
	t.Cleanup(func() { _, _ = request(t, "DELETE", "/queues/test") })

	statusCode, response = request(t, "GET", "/queues/test")

	checkStatusCode(t, statusCode, http.StatusOK)
	checkResponse(t, response, map[string]any{"items_count": 0.0})
}

func TestServer_listHandler(t *testing.T) {
	statusCode, response := request(t, "GET", "/queues")

	checkStatusCode(t, statusCode, http.StatusOK)
	checkResponse(t, response, map[string]any{"queues_info": map[string]any{}, "queues_count": 0.0})

	_, _ = request(t, "POST", "/queues/test")
	t.Cleanup(func() { _, _ = request(t, "DELETE", "/queues/test") })

	statusCode, response = request(t, "GET", "/queues")

	checkStatusCode(t, statusCode, http.StatusOK)

	expectedResponse := map[string]any{
		"queues_info": map[string]any{
			"test": map[string]any{
				"items_count": 0.0,
			},
		},
		"queues_count": 1.0,
	}

	checkResponse(t, response, expectedResponse)
}

func TestServer_deleteHandler(t *testing.T) {
	statusCode, response := request(t, "DELETE", "/queues/test")

	checkStatusCode(t, statusCode, http.StatusNotFound)
	checkResponse(t, response, map[string]any{"message": "a queue with this key does not exist"})

	_, _ = request(t, "POST", "/queues/test")

	statusCode, response = request(t, "DELETE", "/queues/test")

	checkStatusCode(t, statusCode, http.StatusOK)
	checkResponse(t, response, map[string]any{"message": "the queue was successfully deleted"})
}

func request(t *testing.T, method string, url string) (int, map[string]any) {
	t.Helper()

	var res *http.Response
	var err error

	switch method {
	case "POST":
		res, err = http.Post(baseURL+url, "application/json", nil)
	case "DELETE":
		var req *http.Request

		req, err = http.NewRequest(http.MethodDelete, baseURL+url, nil)

		if err != nil {
			t.Fatalf("error creating request, %v", err)
		}

		res, err = http.DefaultClient.Do(req)
	default:
		res, err = http.Get(baseURL + url)
	}

	if err != nil {
		t.Fatalf("request execution error, %v", err)
	}

	var body map[string]any

	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("body decoding error, %v", err)
	}

	if err = res.Body.Close(); err != nil {
		t.Fatalf("body closing error, %v", err)
	}

	return res.StatusCode, body
}

func checkStatusCode(t *testing.T, statusCode int, expectedStatusCode int) {
	t.Helper()

	if statusCode != expectedStatusCode {
		t.Errorf("expected status code %d, received %d", expectedStatusCode, statusCode)
	}
}

func checkResponse(t *testing.T, response map[string]any, expectedResponse map[string]any) {
	t.Helper()

	if !reflect.DeepEqual(expectedResponse, response) {
		t.Errorf("expected response %v, received response %v", expectedResponse, response)
	}
}
