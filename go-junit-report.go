package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wheniwork/go-junit-report/formatter"
	"github.com/wheniwork/go-junit-report/parser"
)

var (
	parserType    = flag.String("parser", "text", "parser to use for the output of go test")
	noXMLHeader   = flag.Bool("no-xml-header", false, "do not print xml header")
	packageName   = flag.String("package-name", "", "specify a package name (compiled test have no package name in output)")
	goVersionFlag = flag.String("go-version", "", "specify the value to use for the go.version property in the generated XML")
	setExitCode   = flag.Bool("set-exit-code", false, "set exit code to 1 if tests failed")
)

func main() {
	flag.Parse()

	if flag.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "%s does not accept positional arguments\n", os.Args[0])
		flag.Usage()
		os.Exit(1)
	}

	var testParser parser.Parser

	switch strings.ToLower(*parserType) {
	case "text":
		testParser = &parser.TextParser{}
	case "json":
		testParser = &parser.JsonParser{}
	default:
		log.Fatalf("parser '%s' is not valid", *parserType)
		return
	}

	// Read input
	report, err := testParser.Parse(os.Stdin, *packageName)
	if err != nil {
		fmt.Printf("Error reading input: %s\n", err)
		os.Exit(1)
	}

	// Write xml
	err = formatter.JUnitReportXML(report, *noXMLHeader, *goVersionFlag, os.Stdout)
	if err != nil {
		fmt.Printf("Error writing XML: %s\n", err)
		os.Exit(1)
	}

	if *setExitCode && report.Failures() > 0 {
		os.Exit(1)
	}
}
