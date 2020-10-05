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
	"net/http"
)

// Route is a specification and  handler for a REST endpoint.
type Route struct {
	verb string
	path string
	fn   func(http.ResponseWriter, *http.Request)
}

// Routes returns a slice of routes of API endpoints to be registered with
// http server
func (g *APIServer) Routes() []*Route {
	h := NewAPIHandler()
	routes := []*Route{
		{verb: "GET", path: "/health", fn: h.Health},
		{verb: "POST", path: versionedPath("/{iac}/{iacVersion}/{cloud}/local/file/scan"), fn: h.scanFile},
		{verb: "POST", path: versionedPath("/{iac}/{iacVersion}/{cloud}/remote/dir/scan"), fn: h.scanRemoteRepo},
	}
	return routes
}

func versionedPath(route string) string {
	return "/" + APIVersion + route
}
