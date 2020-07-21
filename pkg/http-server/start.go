package httpServer

import (
	"log"
	"net/http"
)

func Start() {

	log.Printf("terrascan server listening at port 9010")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":9010", nil))
}
