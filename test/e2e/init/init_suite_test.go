package init_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Init Suite")
}
