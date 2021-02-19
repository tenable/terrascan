// +build tools

package tools

import (
	// used only for testing
	_ "github.com/onsi/ginkgo/ginkgo"
	// used only for testing
	_ "github.com/onsi/gomega"
)

// This file imports packages that are used when running go generate, or used
// during the development process but not otherwise depended on by built code.
