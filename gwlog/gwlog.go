package gwlog

import (
	"encoding/csv"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/user"
	"path"
	"time"
)

type row struct {
	oplog_id      string
	start_date    time.Time
	end_date      time.Time
	source_ip     string
	dest_ip       string
	tool          string
	user_context  string
	command       string
	description   string
	output        string
	comments      string
	operator_name string
}

func (r row) generate() (rowslice []string) {
	rowslice = append(rowslice, r.oplog_id)
	rowslice = append(rowslice, r.start_date.Format("2006-01-02 15:04:05"))
	rowslice = append(rowslice, r.end_date.Format("2006-01-02 15:04:05"))
	rowslice = append(rowslice, r.source_ip)
	rowslice = append(rowslice, r.dest_ip)
	rowslice = append(rowslice, r.tool)
	rowslice = append(rowslice, r.user_context)
	rowslice = append(rowslice, r.command)
	rowslice = append(rowslice, r.description)
	rowslice = append(rowslice, r.output)
	rowslice = append(rowslice, r.comments)
	rowslice = append(rowslice, r.operator_name)
	return
}
func Logwrite(tool, command, description, output string, filename string) {
	// Prepare CSV writer
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error while getting current username: %v", err))
	}
	homeDir := currentUser.HomeDir
	fcsv, err := os.OpenFile(path.Join(homeDir, filename), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error while output CSV file open: %v", err))
	}
	writer := csv.NewWriter(fcsv)

	// Generate row
	r := row{}
	r.oplog_id = os.Getenv("OPLOG_ID")
	r.start_date = time.Now()
	r.end_date = time.Now()
	r.tool = tool
	r.output = output
	r.description = description
	r.command = command
	r.operator_name = currentUser.Username

	// Write record
	err = writer.Write(r.generate())
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error while writing CSV record: %v", err))

	}
	writer.Flush()
	err = fcsv.Close()
	if err != nil {
		log.Fatalf("Error while closing CSV file: %v", err)
	}
}
