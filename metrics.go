package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	res, err := template.New("metric").Parse(fmt.Sprintf(`
        <html>
            <body>
                <h1>Welcome, Chirpy Admin</h1>
                <p>Chirpy has been visited %d times!</p>
            </body>
        </html>`, cfg.fileserverHits.Load()))

	if err != nil {
		log.Fatal(err)
	}
	res.Execute(w, nil)
}
