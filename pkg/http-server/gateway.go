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

// APIGateway implements all the API endpoints, APIGateway should store all the
// metadata info which may be required by all the API handlers
type APIGateway struct{}

// NewAPIGateway returns a new APIGateway{}
func NewAPIGateway() *APIGateway {
	return &APIGateway{}
}
