package main

import (
	"fmt"
	"os"

	"auditexport/internal/cli"
	"auditexport/internal/collectors/github"
	"auditexport/internal/collectors/github/cicd"
	"auditexport/internal/output"
	"auditexport/internal/run"
	"auditexport/internal/summaries"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage: auditexport <command>")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "run":
		run.InitRunContext()

		// --------------------------------------------------
		// Parse & validate flags
		// --------------------------------------------------
		flags := cli.ParseRunFlags(os.Args[2:])

		if flags.Standard != "iso27001" && flags.Standard != "soc2" {
			fmt.Println("invalid --standard value (use iso27001 or soc2)")
			os.Exit(1)
		}

		standard := run.Standard(flags.Standard)
		caps := run.CapabilitiesForStandard(standard)

		// --------------------------------------------------
		// Step count (ISO27001 = 14, SOC2 = 15)
		// --------------------------------------------------
		totalSteps := 14
		if caps.AllowExtendedControls {
			totalSteps = 15
		}
		step := 1

		// --------------------------------------------------
		// 🔐 GitHub Preflight Validation (NO DIRS YET)
		// --------------------------------------------------
		cli.Step(step, totalSteps, "Validating GitHub credentials")
		if err := github.ValidateAuth(); err != nil {
			cli.Failed("GitHub authentication failed")
			fmt.Println(err)
			os.Exit(1)
		}
		cli.Done("GitHub authentication validated")
		step++

		// --------------------------------------------------
		// 🧪 Dry Run Mode (NO FILE SYSTEM / NO API CALLS)
		// --------------------------------------------------
		if flags.DryRun {
			fmt.Println("\nAuditExport — Dry Run Summary")
			fmt.Println("--------------------------------")
			fmt.Printf("Standard        : %s\n", flags.Standard)
			fmt.Printf("Repository      : %s\n", flags.Repo)
			fmt.Printf("Branch          : %s\n", flags.Branch)
			fmt.Printf("Extended Checks : %v\n", caps.AllowExtendedControls)
			fmt.Printf("CI/CD Evidence  : %v\n", caps.AllowCICD)
			fmt.Println("\nNo files were created.")
			fmt.Println("No evidence was collected.")
			fmt.Println("Dry run completed successfully.\n")
			return
		}

		// --------------------------------------------------
		// Initialize run & output directories
		// --------------------------------------------------
		cli.Step(step, totalSteps, "Initializing evidence directory")
		if err := run.InitEvidenceRunDir(); err != nil {
			cli.Failed("failed to initialize evidence directory")
			fmt.Println(err)
			os.Exit(1)
		}
		cli.Done("Evidence directory initialized")
		step++

		cli.Step(step, totalSteps, "Initializing summaries directory")
		if err := run.InitSummariesDir(); err != nil {
			cli.Failed("failed to initialize summaries directory")
			fmt.Println(err)
			os.Exit(1)
		}
		cli.Done("Summaries directory initialized")
		step++

		run.WriteExecutionLog("run started")
		run.WriteExecutionLog("product standard: " + flags.Standard)

		// --------------------------------------------------
		// GitHub Evidence Collection (ISO 27001 baseline)
		// --------------------------------------------------
		cli.Step(step, totalSteps, "Initializing GitHub evidence")
		if err := github.InitGithubDir(); err != nil {
			cli.Failed("failed to initialize GitHub directory")
			fmt.Println(err)
			os.Exit(1)
		}
		cli.Done("GitHub directory initialized")
		step++

		cli.Step(step, totalSteps, "Collecting GitHub organization")
		if err := github.WriteOrganization(); err != nil {
			cli.Failed("failed to collect organization")
			fmt.Println(err)
			os.Exit(1)
		}
		cli.Done("Organization collected")
		step++

		cli.Step(step, totalSteps, "Collecting repositories")
		if err := github.WriteRepositories(); err != nil {
			cli.Failed("failed to collect repositories")
			fmt.Println(err)
			os.Exit(1)
		}
		cli.Done("Repositories collected")
		step++

		cli.Step(step, totalSteps, "Collecting branches")
		if err := github.WriteBranches(); err != nil {
			cli.Failed("failed to collect branches")
			fmt.Println(err)
			os.Exit(1)
		}
		cli.Done("Branches collected")
		step++

		cli.Step(step, totalSteps, "Collecting commits")
		if err := github.WriteCommits(); err != nil {
			cli.Failed("failed to collect commits")
			fmt.Println(err)
			os.Exit(1)
		}
		cli.Done("Commits collected")
		step++

		cli.Step(step, totalSteps, "Collecting pull requests")
		if err := github.WritePullRequests(); err != nil {
			cli.Failed("failed to collect pull requests")
			fmt.Println(err)
			os.Exit(1)
		}
		cli.Done("Pull requests collected")
		step++

		cli.Step(step, totalSteps, "Collecting contributors")
		if err := github.WriteContributors(); err != nil {
			cli.Failed("failed to collect contributors")
			fmt.Println(err)
			os.Exit(1)
		}
		cli.Done("Contributors collected")
		step++

		cli.Step(step, totalSteps, "Collecting access controls")
		if err := github.WriteAccessControls(); err != nil {
			cli.Failed("failed to collect access controls")
			fmt.Println(err)
			os.Exit(1)
		}
		cli.Done("Access controls collected")
		step++

		cli.Step(step, totalSteps, "Collecting protected branches")
		if err := github.WriteProtectedBranches(); err != nil {
			cli.Failed("failed to collect protected branches")
			fmt.Println(err)
			os.Exit(1)
		}
		cli.Done("Protected branches collected")
		step++

		cli.Step(step, totalSteps, "Collecting code owners")
		if err := github.WriteCodeOwners(); err != nil {
			cli.Failed("failed to collect code owners")
			fmt.Println(err)
			os.Exit(1)
		}
		cli.Done("Code owners collected")
		step++

		// --------------------------------------------------
		// SOC 2 Extended Controls (GATED)
		// --------------------------------------------------
		if caps.AllowExtendedControls {

			owner, err := github.GetOwnerFromEvidence()
			if err != nil {
				fmt.Println("failed to resolve github owner:", err)
				os.Exit(1)
			}

			repo := flags.Repo
			branch := flags.Branch

			cli.Step(step, totalSteps, "Collecting CI/CD evidence (GitHub Actions)")
			run.WriteExecutionLog("CI/CD evidence collection enabled")

			if err := cicd.WriteWorkflows(owner, repo); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err := cicd.WriteWorkflowFiles(owner, repo); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err := cicd.WriteWorkflowRuns(owner, repo); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err := cicd.WriteCISummary(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			cli.Done("CI/CD evidence collected")
			step++

			if err := github.WriteRequiredReviews(owner, repo, branch); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err := github.WriteMergePolicies(owner, repo, branch); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// --------------------------------------------------
		// Summaries, Hashes & Packaging
		// --------------------------------------------------
		cli.Step(step, totalSteps, "Generating summaries and packaging evidence")

		if err := run.WriteRunMetadata(flags.Standard); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		_ = summaries.WriteExecutiveSummary()
		_ = summaries.WriteGitHubSummary(flags.Standard)
		_ = summaries.WriteTechnicalSummary()
		_ = summaries.WriteAuditorNotes()

		_ = run.WriteHashes()
		_ = output.ZipEvidence()

		cli.Done("Evidence packaged")
		fmt.Println("\nauditexport run completed")

	default:
		fmt.Println("unknown command:", os.Args[1])
		os.Exit(1)
	}
}
