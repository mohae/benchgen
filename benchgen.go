package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const mainFile = "main.go"

func main() {
	os.Exit(realMain())
}

func realMain() int {
	// this doesn't do anything now, it's mainly to filter out anything
	// that may have been passed as a flag
	flag.Parse()
	// the first arg is the target
	args := flag.Args()

	if len(args) == 0 || args[0] == "" {
		fmt.Println("The path for the new repository must be specified.")
		return 1
	}
	// get GOPATH
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		fmt.Println("Unable to determine the GOPATH.")
		return 1
	}
	// Create the target directory
	dir := filepath.Join(gopath, "src", args[0])
	// See if the path exists, if it does don't do anything
	fi, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println(err)
			return 1
		}
	}
	// if the target exists, dir/main.go must not exist.  The directory may
	// already exist for reasons; e.g. a new repo cloned from github.
	if fi != nil {
		fi, err = os.Stat(filepath.Join(dir, mainFile))
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Println(err)
				return 1
			}
			goto mkdir
		}
		fmt.Printf("%q must not exist: %s\n", mainFile, filepath.Join(dir, mainFile))
		return 1
	}
mkdir:
	// Make the path
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	// write out main.go
	f, err := os.OpenFile(filepath.Join(dir, mainFile), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer f.Close()
	// should probably check for a short write
	_, err = f.WriteString(mainTpl)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}
