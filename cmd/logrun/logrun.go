package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/Ramzeth/logwriter/gwlog"
	"os/exec"
	"strings"
	log "github.com/sirupsen/logrus"
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

	tool := flag.Args()[0]
	// search for tool path
	toolPath, err := exec.LookPath(tool)
	if err != nil {
		log.Fatalf("Unable to find tool: %v",tool)
	}
	toolArgs := flag.Args()[1:]
	cmd := exec.Command(toolPath, toolArgs...)
	output, err := cmd.Output()
	if err != nil {
		log.Errorf(err.Error())
	}
	fmt.Print(string(output))
	description := strings.Join(flag.Args()," ")
	gwlog.Logwrite(tool,description,string(output))

}

