package testsuite_test

import (
	. "github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestTestsuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testsuite Suite")
}
