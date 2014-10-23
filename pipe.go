package main

import (
	"io"
	"sync"
)

type errorChannel chan (error)

func copyStream(to, from io.ReadWriter, wg *sync.WaitGroup, errors errorChannel) {
	err := PipeStream(to, from, wg)
	if err != nil {
		errors <- err
	}
}

func waitForStreams(wg *sync.WaitGroup, errors errorChannel) {
	wg.Wait()
	errors <- nil
}

func PipeStream(to, from io.ReadWriter, wg *sync.WaitGroup) error {
	_, err := io.Copy(to, from)
	wg.Done()
	return err
}

func PipeStreams(a, b io.ReadWriter) error {
	errors := make(errorChannel)
	var wg sync.WaitGroup
	wg.Add(2)
	go copyStream(a, b, &wg, errors)
	go copyStream(b, a, &wg, errors)
	go waitForStreams(&wg, errors)
	return <-errors
}
