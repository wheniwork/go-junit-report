package parser

import (
	"io"
	"regexp"
	"time"
)

// Result represents a test result.
type Result int

// Test result constants
const (
	PASS Result = iota
	FAIL
	SKIP
)

type Parser interface {
	Parse(r io.Reader, pkgName string) (*Report, error)
}

type (
	// Report is a collection of package tests.
	Report struct {
		Packages []Package
	}

	// Package contains the test results of a single package.
	Package struct {
		Name        string
		Duration    time.Duration
		Tests       []*Test
		Benchmarks  []*Benchmark
		CoveragePct string

		// Time is deprecated, use Duration instead.
		Time int // in milliseconds
	}

	// Test contains the results of a single test.
	Test struct {
		Name     string
		Duration time.Duration
		Result   Result
		Output   []string

		SubtestIndent string

		// Time is deprecated, use Duration instead.
		Time int // in milliseconds
	}

	// Benchmark contains the results of a single benchmark.
	Benchmark struct {
		Name     string
		Duration time.Duration
		// number of B/op
		Bytes int
		// number of allocs/op
		Allocs int
	}
)

var (
	regexStatus   = regexp.MustCompile(`--- (PASS|FAIL|SKIP): (.+) \((\d+\.\d+)(?: seconds|s)\)`)
	regexIndent   = regexp.MustCompile(`^([ \t]+)---`)
	regexCoverage = regexp.MustCompile(`^coverage:\s+(\d+\.\d+)%\s+of\s+statements(?:\sin\s.+)?$`)
	regexResult   = regexp.MustCompile(`^(ok|FAIL)\s+([^ ]+)\s+(?:(\d+\.\d+)s|\(cached\)|(\[\w+ failed]))(?:\s+coverage:\s+(\d+\.\d+)%\sof\sstatements(?:\sin\s.+)?)?$`)
	// regexBenchmark captures 3-5 groups: benchmark name, number of times ran, ns/op (with or without decimal), B/op (optional), and allocs/op (optional).
	regexBenchmark       = regexp.MustCompile(`^(Benchmark[^ -]+)(?:-\d+\s+|\s+)(\d+)\s+(\d+|\d+\.\d+)\sns/op(?:\s+(\d+)\sB/op)?(?:\s+(\d+)\sallocs/op)?`)
	regexOutput          = regexp.MustCompile(`(    )*\t(.*)`)
	regexSummary         = regexp.MustCompile(`^(PASS|FAIL|SKIP)$`)
	regexPackageWithTest = regexp.MustCompile(`^# ([^\[\]]+) \[[^\]]+\]$`)
)

// Failures counts the number of failed tests in this report
func (r *Report) Failures() int {
	count := 0

	for _, p := range r.Packages {
		for _, t := range p.Tests {
			if t.Result == FAIL {
				count++
			}
		}
	}

	return count
}
