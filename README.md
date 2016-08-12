# benchgen
Creates benchmark application skeletons that use `github.com\mohae\benchutil`.  This creates the directory path, if it doesn't exist, and the `main.go` file with all of the standard `[benchutil](https://github.com/mohae/benchutil)` flags defined and the basic logic for processing the flags and configuring the benchmark behavior.

This reduces the busy work involved with creating a benchmark application.

This is of no use for situations where using `Benchmark*` testing funcs and running `go test -bench=.` are sufficient.

## Usage

    go install

	benchgen github.com/mohae/benchmark

The above command will create $GOPATH/src/github.com/mohae/benchmark , if it doesn't already exist, and $GOPATH/src/github.com/mohae/benchmark/main.go with the following contents:

```
package main

import (
	"flag"

	"github.com/mohae/benchutil"
)

// flags
var (
	output         string
	format         string
	nameSections   bool
	section        bool
	sectionHeaders bool
	systemInfo     bool
)

func init() {
	flag.StringVar(&output, "output", "stdout", "output destination")
	flag.StringVar(&output, "o", "stdout", "output destination (short)")
	flag.StringVar(&format, "format", "txt", "format of output")
	flag.StringVar(&format, "f", "txt", "format of output")
	flag.BoolVar(&nameSections, "namesections", false, "use group as section name: some restrictions apply")
	flag.BoolVar(&nameSections, "n", false, "use group as section name: some restrictions apply")
	flag.BoolVar(&section, "sections", false, "don't separate groups of tests into sections")
	flag.BoolVar(&section, "s", false, "don't separate groups of tests into sections")
	flag.BoolVar(&sectionHeaders, "sectionheader", false, "if there are sections, add a section header row")
	flag.BoolVar(&sectionHeaders, "h", false, "if there are sections, add a section header row")
	flag.BoolVar(&systemInfo, "sysinfo", false, "add the system information to the output")
	flag.BoolVar(&systemInfo, "i", false, "add the system information to the output")
}

func main() {
	flag.Parse()

	// set up the ticker
	done := make(chan struct{})
	go benchutil.Dot(done)

	// set the output
	var w io.Writer
	var err error
	switch output {
	case "stdout":
		w = os.Stdout
	default:
		w, err = os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer w.(*os.File).Close()
	}
	// get the benchmark for the desired format
	// process the output
	var bench benchutil.Benchmarker
	switch format {
	case "csv":
		bench = benchutil.NewCSVBench(w)
	case "md":
		bench = benchutil.NewMDBench(w)
	default:
		bench = benchutil.NewStringBench(w)
	}
	bench.SectionPerGroup(section)
	bench.SectionHeaders(sectionHeaders)
	bench.IncludeSystemInfo(systemInfo)
	bench.NameSections(nameSections)

	// override column headers (if applicable)

	// run the benchmarks and append the results
	bench.Append(dummyBenchmark())

	// create the output
	fmt.Println("")
	fmt.Println("generating output...")
	err = bench.Out()
	if err != nil {
		fmt.Printf("error generating output: %s\n", err)
	}
}
```

If the output is to have custom column headers, those will need to be set.

All benchmark functions need to be written, called, and appended to the benchmark results.  The `main.go` has a dummy example.

## License
Copyright © 2016: Joel Scoble.  MIT Licensed.  See included LICENSE file for details.
