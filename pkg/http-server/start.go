package httpServer

import (
	"net/http"

	"go.uber.org/zap"
)

func Start() {

	zap.S().Info("terrascan server listening at port 9010")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	zap.S().Fatal(http.ListenAndServe(":9010", nil))
}
