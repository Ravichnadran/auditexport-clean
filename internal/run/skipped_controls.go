package run

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type SkippedControl struct {
	Control  string `json:"control"`
	Reason   string `json:"reason"`
	Details  string `json:"details"`
	Severity string `json:"severity"`
	Impact   string `json:"impact"`
}

type SkippedControlsReport struct {
	CollectedAt time.Time        `json:"collected_at"`
	Standard    string           `json:"standard"`
	Skipped     []SkippedControl `json:"skipped"`
}

var (
	skippedMu sync.Mutex
	skipped   []SkippedControl
)

func RecordSkippedControl(sc SkippedControl) {
	skippedMu.Lock()
	defer skippedMu.Unlock()
	skipped = append(skipped, sc)
}

func WriteSkippedControls(standard string) error {
	if len(skipped) == 0 {
		return nil // nothing to write
	}

	report := SkippedControlsReport{
		CollectedAt: time.Now().UTC(),
		Standard:    standard,
		Skipped:     skipped,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		EvidencePath("run", "skipped_controls.json"),
		data,
		0644,
	)
}
