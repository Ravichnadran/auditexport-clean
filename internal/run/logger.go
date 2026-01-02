package run

import (
	"os"
	"time"
)

func WriteExecutionLog(message string) error {
	entry := time.Now().UTC().Format(time.RFC3339) + " " + message + "\n"

	path := EvidencePath("run", "execution_log.txt")

	f, err := os.OpenFile(	
		path,   
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(entry)
	return err
}
