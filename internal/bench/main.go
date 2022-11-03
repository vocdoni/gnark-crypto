package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type entry struct {
	entry fs.DirEntry
	path  string
}

func main() {
	// quick and dirty helper to benchmark field elements accross branches

	var entries []entry
	err := filepath.WalkDir("../../ecc", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if d.Name() == "fp" || d.Name() == "fr" {
				entries = append(entries, entry{entry: d, path: path})
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	const benchCount = 10
	const regexp = "BBB"
	const refBranch = "developt"
	const newBranch = "feat-addchain"

	benchFileName := "BBB." + runtime.Version() + ".txt"

	var buf bytes.Buffer
	runBenches := func() {
		// checkout(branch)
		for _, e := range entries {
			buf.Reset()
			count := strconv.Itoa(benchCount)
			cmd := exec.Command("go", "test", "-timeout", "10m", "-run", "^$", "-bench", regexp, "-count", count)
			args := strings.Join(cmd.Args, " ")
			log.Println("running benchmark", "dir", e.path, "cmd", args)
			cmd.Dir = e.path
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
			cmd.Stdout = &buf
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}

			if err := os.WriteFile(filepath.Join(e.path, benchFileName), buf.Bytes(), 0600); err != nil {
				log.Fatal(err)
			}
		}
	}

	runBenches()

	// runBenches(refBranch)
	// runBenches(newBranch)

	for _, e := range entries {
		fmt.Println()
		log.Println("benchstat", e.path, regexp)
		cmd := exec.Command("benchstat", benchFileName)
		cmd.Dir = e.path
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println()
	}
}

func checkout(branch string) {
	cmd := exec.Command("git", "checkout", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
