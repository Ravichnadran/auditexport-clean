package main

// --- unchanged imports ---
import (
	"errors"
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

		// --------------------------------------------------
		// 🔒 HARD STOP — FLAG VALIDATION FIRST
		// --------------------------------------------------
		flags := cli.ParseRunFlags(os.Args[2:])

		standard := run.Standard(flags.Standard)
		caps := run.CapabilitiesForStandard(standard)

		// --------------------------------------------------
		// Step count
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
		// 🧪 Dry Run Mode
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
		// Initialize run context & directories
		// --------------------------------------------------
		run.InitRunContext()

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
		// GitHub Evidence Collection
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
			if errors.Is(err, github.ErrAccessControlsPermissionDenied) {
				cli.Skipped("Access controls (insufficient permissions)")

				run.RecordSkippedControl(run.SkippedControl{
					Control:  "github.cicd",
					Reason:   "not_available",
					Details:  "GitHub Actions not enabled or inaccessible",
					Severity: "not_applicable",
					Impact:   "CI/CD controls are not applicable for this repository",
				})

				run.WriteExecutionLog("access controls skipped: insufficient permissions")
			} else {
				cli.Failed("failed to collect access controls")
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			cli.Done("Access controls collected")
		}
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
		// SOC 2 Extended Controls
		// --------------------------------------------------
		if caps.AllowExtendedControls {

			owner, err := github.GetOwnerFromEvidence()
			if err != nil {
				fmt.Println("failed to resolve github owner:", err)
				os.Exit(1)
			}

			repo := flags.Repo
			branch := flags.Branch

			// cli.Step(step, totalSteps, "Collecting CI/CD evidence (GitHub Actions)")
			// if err := cicd.WriteWorkflows(owner, repo); err != nil {
			// 	fmt.Println(err)
			// 	os.Exit(1)
			// }
			// if err := cicd.WriteWorkflowFiles(owner, repo); err != nil {
			// 	fmt.Println(err)
			// 	os.Exit(1)
			// }
			// if err := cicd.WriteWorkflowRuns(owner, repo); err != nil {
			// 	fmt.Println(err)
			// 	os.Exit(1)
			// }
			// if err := cicd.WriteCISummary(); err != nil {
			// 	fmt.Println(err)
			// 	os.Exit(1)
			// }
			// cli.Done("CI/CD evidence collected")
			// step++
			cli.Step(step, totalSteps, "Collecting CI/CD evidence (GitHub Actions)")
			if err := cicd.WriteWorkflows(owner, repo); err != nil {
				if errors.Is(err, cicd.ErrCICDNotAvailable) {
					cli.Skipped("CI/CD evidence (not available)")
					run.WriteExecutionLog("ci/cd skipped: not available")
					step++
				} else {
					cli.Failed("failed to collect CI/CD evidence")
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				_ = cicd.WriteWorkflowFiles(owner, repo)
				_ = cicd.WriteWorkflowRuns(owner, repo)
				_ = cicd.WriteCISummary()
				cli.Done("CI/CD evidence collected")
				step++
			}

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
		_ = run.WriteSkippedControls(flags.Standard)
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
