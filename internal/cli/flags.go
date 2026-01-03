package cli

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type RunFlags struct {
	Standard string
	Repo     string
	Branch   string
	DryRun   bool
	FromDate string
	ToDate   string
}

func ParseRunFlags(args []string) RunFlags {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)

	// Silence default flag output
	fs.SetOutput(os.Stdout)

	standard := fs.String("standard", "", "")
	repo := fs.String("repo", "auditexport", "")
	branch := fs.String("branch", "main", "")
	dryRun := fs.Bool("dry-run", false, "")
	fromDate := fs.String("from-date", "", "")
	toDate := fs.String("to-date", "", "")

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
      GitHub repository name (default: auditexport)

  --branch <branch>
      Target branch name (default: main)

  --from-date YYYY-MM-DD
      Audit window start date (SOC 2 only)

  --to-date YYYY-MM-DD
      Audit window end date (SOC 2 only)

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
  auditexport run --standard soc2 --from-date 2025-10-01 --to-date 2025-12-31
  auditexport run --standard soc2 --dry-run
`)
	}

	// Parse flags
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		os.Exit(2)
	}

	// Reject positional arguments
	if fs.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "unexpected arguments: %v\n\n", fs.Args())
		fs.Usage()
		os.Exit(2)
	}

	// Validate required flag
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

	// Validate date formats (if provided)
	if *fromDate != "" {
		if _, err := time.Parse("2006-01-02", *fromDate); err != nil {
			fmt.Fprintf(os.Stderr, "invalid --from-date format (expected YYYY-MM-DD)\n\n")
			os.Exit(2)
		}
	}

	if *toDate != "" {
		if _, err := time.Parse("2006-01-02", *toDate); err != nil {
			fmt.Fprintf(os.Stderr, "invalid --to-date format (expected YYYY-MM-DD)\n\n")
			os.Exit(2)
		}
	}

	// Logical date validation
	if *fromDate != "" && *toDate != "" {
		from, _ := time.Parse("2006-01-02", *fromDate)
		to, _ := time.Parse("2006-01-02", *toDate)

		if from.After(to) {
			fmt.Fprintln(os.Stderr, "--from-date cannot be after --to-date\n")
			os.Exit(2)
		}
	}

	return RunFlags{
		Standard: *standard,
		Repo:     *repo,
		Branch:   *branch,
		DryRun:   *dryRun,
		FromDate: *fromDate,
		ToDate:   *toDate,
	}
}
