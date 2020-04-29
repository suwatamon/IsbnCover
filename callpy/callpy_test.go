package callpy

import (
	"sync"
	"testing"
)

func TestWithChan(t *testing.T) {
	var (
		chIn  = make(chan string, 1)
		chOut = make(chan string, 1)
		mu    sync.Mutex
	)
	go WithChan("./test.py", chIn, chOut)

	inputSample := "test string"
	mu.Lock()
	chIn <- inputSample
	outputSample := <-chOut
	mu.Unlock()

	if inputSample != outputSample {
		t.Errorf("input and output string are differ: %s, %s\n", inputSample, outputSample)
	}
	close(chIn)
	close(chOut)
}
