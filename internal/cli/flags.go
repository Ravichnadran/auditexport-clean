package cli

import (
	"flag"
	"fmt"
	"os"
)

type RunFlags struct {
	Standard string
	Repo     string
	Branch   string
	DryRun   bool
}

func ParseRunFlags(args []string) RunFlags {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)

	// Silence default Go error output
	fs.SetOutput(os.Stdout)

	// -------------------------------
	// Flags (ALL inputs are flags)
	// -------------------------------
	standard := fs.String("standard", "", "")
	repo := fs.String("repo", "auditexport", "")
	branch := fs.String("branch", "main", "")
	dryRun := fs.Bool("dry-run", false, "")

	// -------------------------------
	// Custom professional help output
	// -------------------------------
	fs.Usage = func() {
		fmt.Fprintln(os.Stdout, `
AuditExport — Technical Audit Evidence Generator

USAGE:
  auditexport run [flags]

REQUIRED:
  --standard iso27001|soc2
      Compliance standard to generate evidence for

OPTIONAL:
  --repo <repository>
      GitHub repository name
      Default: auditexport

  --branch <branch>
      Target branch name
      Default: main

  --dry-run
      Validate configuration and print execution plan
      without creating files or collecting evidence

  --help
      Show this help message and exit

ENVIRONMENT:
  GITHUB_TOKEN
      Required for GitHub evidence collection (read-only access)

EXAMPLES:
  auditexport run --standard iso27001
  auditexport run --standard soc2 --repo my-repo
  auditexport run --standard soc2 --dry-run
`)
	}

	// -------------------------------
	// Parse flags
	// -------------------------------
	if err := fs.Parse(args); err != nil {
		// Handles --help and unknown flags cleanly
		os.Exit(0)
	}

	// -------------------------------
	// Reject positional arguments
	// -------------------------------
	if fs.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "unexpected arguments: %v\n\n", fs.Args())
		fs.Usage()
		os.Exit(2)
	}

	// -------------------------------
	// Validate required flags
	// -------------------------------
	if *standard == "" {
		fmt.Fprintln(os.Stderr, "error: --standard is required\n")
		fs.Usage()
		os.Exit(2)
	}

	if *standard != "iso27001" && *standard != "soc2" {
		fmt.Fprintf(os.Stderr, "invalid --standard value: %s\n\n", *standard)
		fs.Usage()
		os.Exit(2)
	}

	return RunFlags{
		Standard: *standard,
		Repo:     *repo,
		Branch:   *branch,
		DryRun:   *dryRun,
	}
}
