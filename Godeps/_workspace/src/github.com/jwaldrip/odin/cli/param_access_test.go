package cli_test

import (
	. "github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/jwaldrip/odin/cli"

	. "github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/onsi/gomega"
)

var _ = Describe("Param Access", func() {

	var cli *CLI
	var cmd Command
	var didRun bool

	BeforeEach(func() {
		didRun = false
		runFn := func(c Command) {
			cmd = c
			didRun = true
		}
		cli = New("v1.0.0", "sample description", runFn, "foo", "bar")
		cli.ErrorHandling = PanicOnError
		cli.Mute()
		cli.Start("cmd", "a", "b")
	})

	Describe("Flag(string) Value", func() {
		It("should return a flag's value", func() {
			Expect(cmd.Param("foo").Get()).To(Equal("a"))
		})

		Context("when a param is not defined", func() {
			It("should panic", func() {
				Expect(func() { cmd.Param("undefined") }).Should(Panic())
			})
		})
	})

	Describe("Flags() map[string]value", func() {
		It("should return each flags value", func() {
			for k, v := range cmd.Params() {
				Expect(v).To(Equal(cmd.Param(k)))
			}
		})
	})

})
