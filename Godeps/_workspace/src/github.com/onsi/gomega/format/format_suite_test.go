package format_test

import (
	. "github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestFormat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Format Suite")
}
