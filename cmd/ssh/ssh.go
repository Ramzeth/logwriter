package main

import (
	"flag"
	"github.com/Ramzeth/logwriter/gwlog"
	"github.com/creack/pty"
	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	flag.Parse()
	tool := "ssh"
	// search for tool path
	toolPath, err := exec.LookPath(tool)
	if err != nil {
		log.Fatalf("Unable to find ssh: %v", tool)
	}

	var toolArgs = flag.Args()
	user_context := ""
	dest_ip := ""
	for _, arg := range toolArgs {
		// get username and host
		if strings.Contains(arg, "@") {
			ssh_uri_tokens := strings.Split(arg, "@")
			user_context = ssh_uri_tokens[0]
			dest_ip = ssh_uri_tokens[1]
		}
	}

	// ToDo implement source_ip gather
	source_ip := ""
	command := strings.Join(flag.Args(), " ")
	gwlog.Logwrite(source_ip, dest_ip, tool, user_context, command, "ssh log wrapper", "", "", "gwlog.csv")

	// From https://github.com/creack/pty
	cmd := exec.Command(toolPath, toolArgs...)
	ptmx, err := pty.Start(cmd)
	defer func() { _ = ptmx.Close() }()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Errorf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH
	defer func() { signal.Stop(ch); close(ch) }()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)
}
