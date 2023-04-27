package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/html")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(fmt.Sprintf(`
<html>
	<body>
<h1>Welcome, Chirpy Admin</h1>
<p>Chirpy has been visited %d times!</p>
	</body>
	</html>

`, cfg.fileserverHits)))

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cfg.fileserverHits++

		next.ServeHTTP(writer, request)
	})

}
