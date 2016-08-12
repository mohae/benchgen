package main

var mainTpl = `package main

import (
	"flag"

	"github.com/mohae/benchutil"
)

// flags
var (
	output         string
	format         string
	help           bool
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
	flag.BoolVar(&help, "help", false, "help")
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&nameSections, "namesections", false, "use group as section name: some restrictions apply")
	flag.BoolVar(&nameSections, "n", false, "use group as section name: some restrictions apply")
	flag.BoolVar(&section, "sections", false, "don't separate groups of tests into sections")
	flag.BoolVar(&section, "s", false, "don't separate groups of tests into sections")
	flag.BoolVar(&sectionHeaders, "sectionheader", false, "if there are sections, add a section header row")
	flag.BoolVar(&sectionHeaders, "r", false, "if there are sections, add a section header row")
	flag.BoolVar(&systemInfo, "sysinfo", false, "add the system information to the output")
	flag.BoolVar(&systemInfo, "i", false, "add the system information to the output")
}

func main() {
	os.Exit(realMain())
}

// The actual work is done here; allows defers to run.
func realMain() int {
	flag.Parse()
	// check the args to see if help was passed w/o a -
	args := flag.Args()
	for _, v := range args {
		if v == "help" {
			help = true
			break
		}
	}
	if help {
		flag.Usage()
		return 1
	}

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
			return(1)
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
	//	bench.SetGroupColumnHeader("")
	//	bench.SetSubGroupColumnHeader("")
	//	bench.SetNameColumnHeader("")
	//	bench.SetDescColumnHeader("")
	//	bench.SetNoteColumnHeader("")

	// run the benchmarks and append the results
	bench.Append(dummyBenchmark())

	// create the output
	fmt.Println("")
	fmt.Println("generating output...")
	err = bench.Out()
	if err != nil {
		fmt.Printf("error generating output: %s\n", err)
		return 1
	}
	return 0
}
`
