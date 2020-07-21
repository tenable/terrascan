package awsProvider

import (
	"log"
)

type AWSProvider struct{}

// CreateNormalizedJson creates a normalized json for the given input
func (a *AWSProvider) CreateNormalizedJson() {
	log.Printf("creating normalized json for AWS resources")
}
