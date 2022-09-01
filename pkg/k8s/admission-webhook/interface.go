/*
    Copyright (C) 2022 Tenable, Inc.

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

package admissionwebhook

import (
	admissionv1 "k8s.io/api/admission/v1"
)

// AdmissionWebhook interface needs to be implemented by all k8s admission
// webhooks i.e validating and mutating webhooks
type AdmissionWebhook interface {

	// Authorize checks if the incoming webhooks have valid apiKey
	Authorize(apiKey string) error

	// DecodeAdmissionReviewRequest reads the incoming admission request body
	// and decodes it into an AdmissionReviewRequest struct
	DecodeAdmissionReviewRequest(payload []byte) (admissionv1.AdmissionReview, error)

	// ProcessWebhook processes the incoming AdmissionReview and creates
	// a AdmissionResponse
	ProcessWebhook(review admissionv1.AdmissionReview, serverURL string) (*admissionv1.AdmissionReview, error)
}
