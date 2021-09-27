package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/Ramzeth/logwriter/gwlog"
	"strings"
)


func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s Description\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main()  {
	// Check flags
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	description := strings.Join(flag.Args()," ")
	gwlog.Logwrite("Custom operator record",description,"")
}

