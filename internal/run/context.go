package run

import (
	"fmt"
	"time"
)

var (
	runTimestamp string
	baseDir      string
)

func InitRunContext() {
	if baseDir != "" {
		return
	}

	runTimestamp = time.Now().UTC().Format("2006-01-02T15-04-05Z")
	baseDir = fmt.Sprintf("auditexport-evidence-%s", runTimestamp)
}

func BaseDir() string {
	InitRunContext()
	return baseDir
}
