package parser

import (
	"bufio"
	"encoding/json"
	"io"
	"time"
)

var (
	_ Parser = &JsonParser{}
)

type JsonParser struct{}

func (j *JsonParser) Parse(r io.Reader, pkgName string) (*Report, error) {
	reader := bufio.NewReader(r)

	tests := map[string]map[string]Test{}

	type TestEvent struct {
		Time    time.Time `json:"time"` // encodes as an RFC3339-format string
		Action  string    `json:"action"`
		Package string    `json:"package"`
		Test    string    `json:"test"`
		Elapsed float64   `json:"elapsed"` // seconds
		Output  string    `json:"output"`
	}

	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break // If we have reached the end of the file or stream then exit the parse loop.
		}

		var item TestEvent
		if err := json.Unmarshal(lineBytes, &item); err != nil {
			panic(err)
		}

		switch len(item.Test) {
		case 0:
		// Not associated with a test.
		default:
			pkg, ok := tests[item.Package]
			if !ok {
				pkg = map[string]Test{}
				tests[item.Package] = pkg
			}

			test, ok := pkg[item.Test]
			if !ok {
				test = Test{
					Name:   item.Test,
					Output: make([]string, 0),
				}
				tests[item.Package][item.Test] = test
			}

			switch item.Action {
			case "output":
				test.Output = append(test.Output, item.Output)
			case "pass":
				test.Duration = time.Duration(item.Elapsed) * time.Second
				test.Result = PASS
			case "fail":
				test.Duration = time.Duration(item.Elapsed) * time.Second
				test.Result = FAIL
			case "skip":
				test.Duration = time.Duration(item.Elapsed) * time.Second
				test.Result = SKIP
			}

			tests[item.Package][item.Test] = test
		}
	}

	report := &Report{
		Packages: make([]Package, 0),
	}

	for pkg, packageTests := range tests {
		testList := make([]*Test, 0, len(packageTests))
		for _, test := range packageTests {
			testList = append(testList, &test)
		}


		report.Packages = append(report.Packages, Package{
			Name:        pkg,
			Tests:       testList,
		})
	}

	return report, nil
}
