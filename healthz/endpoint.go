package endpoint

import (
	"net/http"
)

func EndpointHandler() http.Handler {
	respWriter := http.NewServeMux()
	respWriter.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	return respWriter
}
