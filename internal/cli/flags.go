package cli

import "flag"

type RunFlags struct {
	Standard string
	Repo     string
	Branch   string
}

func ParseRunFlags(args []string) RunFlags {
	fs := flag.NewFlagSet("run", flag.ExitOnError)

	standard := fs.String(
		"standard",
		"iso27001",
		"compliance standard: iso27001 or soc2",
	)

	repo := fs.String(
		"repo",
		"auditexport",
		"target repository name",
	)

	branch := fs.String(
		"branch",
		"main",
		"target branch name",
	)

	// SAFE: parse only the args passed to this function
	_ = fs.Parse(args)

	return RunFlags{
		Standard: *standard,
		Repo:     *repo,
		Branch:   *branch,
	}
}
