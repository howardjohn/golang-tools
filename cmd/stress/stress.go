// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

// The stress utility is intended for catching sporadic failures.
// It runs a given process in parallel in a loop and collects any failures.
// Usage:
// 	$ stress ./fmt.test -test.run=TestSometing -test.cpu=10
// You can also specify a number of parallel processes with -p flag;
// instruct the utility to not kill hanged processes for gdb attach;
// or specify the failure output you are looking for (if you want to
// ignore some other sporadic failures).
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"syscall"
	"time"
)

var (
	flagP        = flag.Int("p", runtime.NumCPU(), "run `N` processes in parallel")
	flagLimit    = flag.Int("limit", 100, "maximum number of failures until exiting")
	flagTimeout  = flag.Duration("timeout", 10*time.Minute, "timeout each process after `duration`")
	flagKill     = flag.Bool("kill", true, "kill timed out processes if true, otherwise just print pid (to attach with gdb)")
	flagFailure  = flag.String("failure", "", "fail only if output matches `regexp`")
	flagIgnore   = flag.String("ignore", "", "ignore failure if output matches `regexp`")
	flagOutput   = flag.String("o", defaultPrefix(), "output failure logs to `path` plus a unique suffix")
	flagFailFast = flag.Bool("f", false, "exit on first failure")
)

func init() {
	flag.Usage = func() {
		os.Stderr.WriteString(`The stress utility is intended for catching sporadic failures.
It runs a given process in parallel in a loop and collects any failures.
Usage:

	$ stress ./fmt.test -test.run=TestSometing -test.cpu=10

`)
		flag.PrintDefaults()
	}
}

func defaultPrefix() string {
	date := time.Now().Format("go-stress-20060102T150405-")
	return filepath.Join(os.TempDir(), date)
}

type result struct {
	out      []byte
	duration time.Duration
}

func main() {
	flag.Parse()
	if *flagP <= 0 || *flagTimeout <= 0 || len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	var failureRe, ignoreRe *regexp.Regexp
	if *flagFailure != "" {
		var err error
		if failureRe, err = regexp.Compile(*flagFailure); err != nil {
			fmt.Println("bad failure regexp:", err)
			os.Exit(1)
		}
	}
	if *flagIgnore != "" {
		var err error
		if ignoreRe, err = regexp.Compile(*flagIgnore); err != nil {
			fmt.Println("bad ignore regexp:", err)
			os.Exit(1)
		}
	}
	results := make(chan result)
	for i := 0; i < *flagP; i++ {
		go func() {
			for {
				t0 := time.Now()
				cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
				done := make(chan bool)
				if *flagTimeout > 0 {
					go func() {
						select {
						case <-done:
							return
						case <-time.After(*flagTimeout):
						}
						if !*flagKill {
							fmt.Printf("process %v timed out\n", cmd.Process.Pid)
							return
						}
						cmd.Process.Signal(syscall.SIGABRT)
						select {
						case <-done:
							return
						case <-time.After(10 * time.Second):
						}
						cmd.Process.Kill()
					}()
				}
				out, err := cmd.CombinedOutput()
				close(done)
				if err != nil && (failureRe == nil || failureRe.Match(out)) && (ignoreRe == nil || !ignoreRe.Match(out)) {
					out = append(out, fmt.Sprintf("\nERROR: %v", err)...)
				} else {
					out = []byte{}
				}
				results <- result{out, time.Since(t0)}
			}
		}()
	}
	runs, fails := 0, 0
	totalDuration := time.Duration(0)
	max, min := time.Duration(0), time.Duration(1<<63-1)
	ticker := time.NewTicker(2 * time.Second).C
	displayProgress := func() {
		total := runs + fails
		if total == 0 {
			fmt.Printf("no runs so far\n")
		} else {
			fmt.Printf("%v runs so far, %v failures (%.2f%% pass rate). %v avg, %v max, %v min\n",
				runs, fails, 100.0*(float64(runs)/float64(total)), totalDuration/time.Duration(total), max, min)
		}
	}
	for {
		select {
		case res := <-results:
			runs++
			totalDuration += res.duration
			if res.duration > max {
				max = res.duration
			}
			if res.duration < min {
				min = res.duration
			}
			if runs == 1 {
				displayProgress()
			}
			if len(res.out) == 0 {
				continue
			}
			fails++
			displayProgress()
			dir, path := filepath.Split(*flagOutput)
			f, err := ioutil.TempFile(dir, path)
			if err != nil {
				fmt.Printf("failed to create temp file: %v\n", err)
				os.Exit(1)
			}
			f.Write(res.out)
			f.Close()
			if len(res.out) > 2<<10 {
				fmt.Printf("%s\n%s\n…\n", f.Name(), res.out[:2<<10])
			} else {
				fmt.Printf("%s\n%s\n", f.Name(), res.out)
			}
			if *flagFailFast {
				fmt.Printf("fail fast enabled, exiting\n")
				os.Exit(1)
			}
			if fails >= *flagLimit {
				fmt.Printf("failure limit hit, exiting\n")
				os.Exit(1)
			}
		case <-ticker:
			displayProgress()
		}
	}
}
