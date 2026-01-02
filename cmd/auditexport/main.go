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

		totalSteps := 15
		step := 1

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

			if caps.AllowCICD {
				cli.Step(step, totalSteps, "Collecting CI/CD evidence (GitHub Actions)")
				run.WriteExecutionLog("CI/CD evidence collection enabled")

				if err := cicd.WriteWorkflows(owner, repo); err != nil {
					cli.Failed("failed to collect workflows")
					fmt.Println(err)
					os.Exit(1)
				}

				if err := cicd.WriteWorkflowFiles(owner, repo); err != nil {
					cli.Failed("failed to collect workflow files")
					fmt.Println(err)
					os.Exit(1)
				}

				if err := cicd.WriteWorkflowRuns(owner, repo); err != nil {
					cli.Failed("failed to collect workflow runs")
					fmt.Println(err)
					os.Exit(1)
				}

				if err := cicd.WriteCISummary(); err != nil {
					cli.Failed("failed to write CI summary")
					fmt.Println(err)
					os.Exit(1)
				}

				cli.Done("CI/CD evidence collected")
				run.WriteExecutionLog("CI/CD evidence collected")
				step++
			} else {
				cli.Skipped("CI/CD evidence (ISO 27001)")
				run.WriteExecutionLog("CI/CD evidence skipped (ISO 27001)")
			}

			if err := github.WriteRequiredReviews(owner, repo, branch); err != nil {
				fmt.Println("failed to write github required reviews:", err)
				os.Exit(1)
			}

			if err := github.WriteMergePolicies(owner, repo, branch); err != nil {
				fmt.Println("failed to write github merge policies:", err)
				os.Exit(1)
			}

			run.WriteExecutionLog("github required reviews collected")
			run.WriteExecutionLog("github merge policies collected")

		} else {
			cli.Skipped("SOC2 extended controls")
			run.WriteExecutionLog("required reviews skipped (ISO 27001)")
			run.WriteExecutionLog("merge policies skipped (ISO 27001)")
		}

		// --------------------------------------------------
		// Run Metadata & Summaries
		// --------------------------------------------------
		cli.Step(step, totalSteps, "Generating summaries and metadata")
		if err := run.WriteRunMetadata(flags.Standard); err != nil {
			fmt.Println("failed to write run metadata:", err)
			os.Exit(1)
		}

		if err := summaries.WriteExecutiveSummary(); err != nil {
			fmt.Println("failed to write executive summary:", err)
			os.Exit(1)
		}

		if err := summaries.WriteGitHubSummary(); err != nil {
			fmt.Println("failed to write github summary:", err)
			os.Exit(1)
		}

		if err := summaries.WriteTechnicalSummary(); err != nil {
			fmt.Println("failed to write technical summary:", err)
			os.Exit(1)
		}

		if err := summaries.WriteAuditorNotes(); err != nil {
			fmt.Println("failed to write auditor notes:", err)
			os.Exit(1)
		}
		cli.Done("Summaries generated")
		step++

		// --------------------------------------------------
		// Integrity & Packaging
		// --------------------------------------------------
		cli.Step(step, totalSteps, "Generating hashes and packaging evidence")
		if err := run.WriteHashes(); err != nil {
			fmt.Println("failed to write hashes:", err)
			os.Exit(1)
		}

		if err := output.ZipEvidence(); err != nil {
			fmt.Println("failed to zip evidence:", err)
			os.Exit(1)
		}
		cli.Done("Evidence packaged")

		run.WriteExecutionLog("hashes generated")
		run.WriteExecutionLog("evidence zipped")
		run.WriteExecutionLog("summaries generated")
		run.WriteExecutionLog("run completed")

		fmt.Println("\nauditexport run completed")

	default:
		fmt.Println("unknown command:", os.Args[1])
		os.Exit(1)
	}
}
