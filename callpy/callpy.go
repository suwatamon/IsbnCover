package callpy

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// WithChan calls Python script; stdio is connected cannels of arguments
func WithChan(pyScript string, chIn <-chan string, chOut chan<- string) {
	execpy := exec.Command("python3", pyScript)
	stdin, err := execpy.StdinPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	stdout, err := execpy.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	stderr, err := execpy.StderrPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	scannerErr := bufio.NewScanner(stderr)
	go func() {
		for scannerErr.Scan() {
			fmt.Fprintln(os.Stderr, scannerErr.Text())
		}
	}()

	execpy.Start()
	defer execpy.Wait()

	go func() {
		for {
			str, ok := <-chIn
			if ok == false {
				stdin.Close()
				return
			}
			io.WriteString(stdin, str+"\n")
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			chOut <- scanner.Text()
		}
	}()
}
