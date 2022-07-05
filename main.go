// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gotest is a tiny program that shells out to `go test`
// and prints the output in color.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

var (
	pass = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	skip = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	fail = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	skipnotest bool
)

const (
	colorsEnv      = "GOTEST_COLORS"
	skipNoTestsEnv = "GOTEST_SKIPNOTESTS"
)

func main() {
	setColors()
	enableSkipNoTests()

	os.Exit(gotest(os.Args[1:]))
}

func gotest(args []string) int {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	r, w := io.Pipe()
	defer w.Close()

	args = append([]string{"test"}, args...)
	cmd := exec.Command("go", args...)
	cmd.Stderr = w
	cmd.Stdout = w
	cmd.Env = os.Environ()

	if err := cmd.Start(); err != nil {
		log.Print(err)
		wg.Done()
		return 1
	}

	go consume(&wg, r)

	sigc := make(chan os.Signal)
	done := make(chan struct{})
	defer func() {
		done <- struct{}{}
	}()
	signal.Notify(sigc)

	go func() {
		for {
			select {
			case sig := <-sigc:
				cmd.Process.Signal(sig)
			case <-done:
				return
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
			return ws.ExitStatus()
		}
		return 1
	}
	return 0
}

func consume(wg *sync.WaitGroup, r io.Reader) {
	defer wg.Done()
	reader := bufio.NewReader(r)
	for {
		l, _, err := reader.ReadLine()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Print(err)
			return
		}
		parse(string(l))
	}
}

func parse(line string) {
	trimmed := strings.TrimSpace(line)
	defer color.Unset()

	var style lipgloss.Style
	switch {
	case strings.Contains(trimmed, "[no test files]"):
		if skipnotest {
			return
		}

    // passed
	case strings.HasPrefix(trimmed, "--- PASS"):
		fallthrough
	case strings.HasPrefix(trimmed, "ok"):
		fallthrough
	case strings.HasPrefix(trimmed, "PASS"):
		style = pass

	// skipped
	case strings.HasPrefix(trimmed, "--- SKIP"):
		style = skip

	// failed
	case strings.HasPrefix(trimmed, "--- FAIL"):
		fallthrough
	case strings.HasPrefix(trimmed, "FAIL"):
		style = fail
	}

	fmt.Printf("%s\n", style.Render(line))
}

func setColors() {
	v := os.Getenv(colorsEnv)
	if v == "" {
		return
	}
	vals := strings.Split(v, ",")
	if len(vals) != 3 {
		return
	}

	if vals[0] != "" {
		pass = fail.Copy().Foreground(lipgloss.Color(vals[0]))
	}

	if vals[1] != "" {
		fail = fail.Copy().Foreground(lipgloss.Color(vals[1]))
	}

	if vals[2] != "" {
		skip = fail.Copy().Foreground(lipgloss.Color(vals[2]))
	}
}

func enableSkipNoTests() {
	v := os.Getenv(skipNoTestsEnv)
	if v == "" {
		return
	}
	v = strings.ToLower(v)
	skipnotest = v == "true"
}
