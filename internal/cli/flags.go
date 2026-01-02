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
}

func ParseRunFlags(args []string) RunFlags {
	fs := flag.NewFlagSet("run", flag.ExitOnError)

	// Custom professional help output
	fs.Usage = func() {
		fmt.Fprintln(os.Stdout, `
AuditExport  —  Technical Audit Evidence Generator

USAGE       :  auditexport run [flags]

REQUIRED    :
  --standard    iso27001|soc2
                Compliance standard to generate evidence for

OPTIONAL    :
  --repo        <name>
                GitHub repository name (default: auditexport)

  --branch      <name>
                Target branch name (default: main)

  --help
                Show this help message and exit

ENVIRONMENT :
  GITHUB_TOKEN
                Required for GitHub evidence collection (read-only)

EXAMPLES    :
  auditexport run --standard iso27001
  auditexport run --standard soc2 --repo my-repo
  auditexport run --help
`)

	}

	standard := fs.String("standard", "iso27001", "")
	repo := fs.String("repo", "auditexport", "")
	branch := fs.String("branch", "main", "")

	_ = fs.Parse(args)

	return RunFlags{
		Standard: *standard,
		Repo:     *repo,
		Branch:   *branch,
	}
}
