package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/Ramzeth/logwriter/gwlog"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

func main() {
	tool := "ssh"
	// search for tool path
	toolPath, err := exec.LookPath(tool)
	if err != nil {
		log.Fatalf("Unable to find ssh: %v", tool)
	}

	toolArgs := flag.Args()[1:]
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

	cmd := exec.Command(toolPath, toolArgs...)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	output := ""
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		output = output + m + "\n"
	}
	cmd.Wait()

}
