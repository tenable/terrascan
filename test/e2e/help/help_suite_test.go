package help_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHelp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Help Suite")
}
