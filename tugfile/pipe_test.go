package tugfile

import (
	"bytes"
	"io"
	"sync"
	"testing"

	. "github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/onsi/gomega"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pipe")
}

type stream struct {
	readBuffer  string
	writeBuffer string
}

func (rw *stream) Read(p []byte) (int, error) {
	var err error
	n := copy(p, rw.readBuffer)
	if n > 0 && n == len(p) {
		err = io.EOF
	}

	return n, err
}

func (rw *stream) Write(p []byte) (int, error) {
	rw.writeBuffer = string(p)
	return len(rw.writeBuffer), io.EOF
}

var _ = Describe("pipe", func() {
	Describe("PipeStream", func() {
		It("should stream data from one io.ReadWriter to the other", func() {
			to := bytes.NewBuffer(nil)
			from := bytes.NewBuffer(nil)

			var wg sync.WaitGroup
			wg.Add(1)

			from.WriteString("echo")
			PipeStream(to, from, &wg)
			wg.Wait()

			Expect(to.String()).To(Equal("echo"))
		})
	})

	Describe("PipeStreams", func() {
		It("should relay traffic from one ReadWriter to the other", func() {
			s1, s2 := &stream{}, &stream{}

			go func() {
				s1.readBuffer = "s1"
			}()
			go func() {
				s2.readBuffer = "s2"
			}()

			PipeStreams(s1, s2)
			// s2's readBuffer is copied over to s1's writeBuffer
			Expect(s1.writeBuffer).To(Equal("s2"))
			// s1's readBuffer is copied over to s2's writeBuffer
			Expect(s2.writeBuffer).To(Equal("s1"))
		})
	})
})
