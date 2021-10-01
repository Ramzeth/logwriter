package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/Ramzeth/logwriter/gwlog"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s Description\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
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
		log.Fatalf("Unable to find tool: %v", tool)
	}
	toolArgs := flag.Args()[1:]
	cmd := exec.Command(toolPath, toolArgs...)
	prout, pwout := io.Pipe()
	prerr, pwerr := io.Pipe()
	cmd.Stdout = pwout
	cmd.Stderr = pwerr

	tout := io.TeeReader(prout, os.Stdout)
	terr := io.TeeReader(prerr, os.Stderr)

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	var bout, berr bytes.Buffer

	go func() {
		if _, err := io.Copy(&bout, tout); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if _, err := io.Copy(&berr, terr); err != nil {
			log.Fatal(err)
		}
	}()

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	command := strings.Join(flag.Args(), " ")
	gwlog.Logwrite("", "", tool, "", command, "Custom tool with logrun", string(bout.String()+berr.String()), "", "gwlog.csv")

}
