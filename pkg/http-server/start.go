/*
    Copyright (C) 2020 Accurics, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package httpserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/accurics/terrascan/pkg/logging"
	"github.com/gorilla/mux"
)

// Start initializes api routes and starts http server
func Start(port, configFile, certFile, privateKeyFile string) {

	// create a new API server
	server := NewAPIServer()

	// get all routes
	routes := server.Routes(configFile)

	// get port
	if os.Getenv(TerrascanServerPort) != "" {
		port = os.Getenv(TerrascanServerPort)
	}

	// register routes and start the http server
	server.start(routes, port, certFile, privateKeyFile)
}

// start http server
func (g *APIServer) start(routes []*Route, port, certFile, privateKeyFile string) {

	var (
		err    error
		logger = logging.GetDefaultLogger() // new logger
		router = mux.NewRouter()            // new router
	)

	logger.Info("registering routes...")

	// register all routes
	for _, v := range routes {
		logger.Info("Route ", v.verb, " - ", v.path)
		router.Methods(v.verb).Path(v.path).HandlerFunc(v.fn)
	}

	// Add a route for all static templates / assets. Currently used for the Webhook logs views
	// go/terrascan/asset is the path where the assets files are located inside the docker container
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("/go/terrascan/assets"))))

	// start http server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		var err error
		if certFile != "" {
			// In case a certificate file is specified, the server support TLS
			err = server.ListenAndServeTLS(certFile, privateKeyFile)
		} else {
			err = server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()
	logger.Infof("http server listening at port %v", port)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// try to stop the server gracefully with default 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		logger.Fatalf("server failed to exit gracefully. error: '%v'", err)
	}
	logger.Info("server exiting gracefully")
}
