package cli

import (
	"fmt"
)

// Step prints a progress step to STDOUT (user-facing only)
func Step(index int, total int, message string) {
	fmt.Printf("[%d/%d] %s...\n", index, total, message)
}

// Done marks a step as completed
func Done(message string) {
	fmt.Printf("      ✓ %s\n", message)
}

// Skipped marks a skipped step
func Skipped(message string) {
	fmt.Printf("      ↷ %s\n", message)
}

// Failed marks a failed step (before exit)
func Failed(message string) {
	fmt.Printf("      ✗ %s\n", message)
}
