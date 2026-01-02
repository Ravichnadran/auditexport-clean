# AuditExport

**AuditExport** is a command-line tool that automatically collects **auditor-ready technical evidence** from GitHub repositories for **ISO/IEC 27001** and **SOC 2** audits.

It converts what is normally **days of manual evidence gathering** into a **single deterministic export**.

---

## What problem does this solve?

During audits, teams struggle with:

* Manually collecting GitHub screenshots
* Explaining branch protection and review rules repeatedly
* Proving access controls and change management
* Re-exporting the same evidence every audit cycle

AuditExport solves this by generating **verifiable, structured, read-only evidence** that auditors already understand.

---

## What AuditExport generates

AuditExport produces a ZIP file containing:

### GitHub baseline evidence (ISO 27001)

* Organization details
* Repository inventory
* Branches
* Commits
* Pull requests
* Contributors
* Access controls
* Protected branches
* CODEOWNERS
* Execution log
* Evidence integrity hashes
* Auditor-friendly summaries

### Additional SOC 2 controls (SOC 2 edition only)

* Required pull-request reviews
* Merge policy enforcement
* CI/CD (GitHub Actions) evidence
* Change-control enforcement proof

All evidence is exported as **JSON + Markdown**, ready for auditors.

---

## Supported standards

| Standard      | Supported             |
| ------------- | --------------------- |
| ISO/IEC 27001 | ✅ Yes                 |
| SOC 2         | ✅ Yes (extended mode) |

---

## Security model (important)

AuditExport is designed to be **safe by default**:

* Read-only GitHub API access
* No secrets stored
* No credentials written to disk
* No artifacts downloaded
* No builds executed
* No CI logs parsed
* No screenshots required

The tool only **observes and exports metadata**.

---

## What AuditExport does NOT do

To keep audits predictable and safe, AuditExport intentionally does **not**:

* Execute CI/CD pipelines
* Download build artifacts
* Read secret values
* Parse application logs
* Modify GitHub settings
* Perform security scanning

This keeps the tool **auditor-friendly and non-intrusive**.

---

## Output structure (example)

```
evidence/
 ├── github/
 │   ├── organization.json
 │   ├── repositories.json
 │   ├── branches.json
 │   ├── access_controls.json
 │   ├── protected_branches.json
 │   ├── code_owners.json
 │   └── cicd/              # SOC 2 only
 ├── run/
 │   ├── run_metadata.json
 │   ├── execution_log.txt
 │   └── hashes.json
 ├── summaries/
 │   ├── executive_summary.md
 │   ├── github_summary.md
 │   └── technical_summary.md
 └── evidence.zip
```

---

## Licensing

This is a **commercial tool**.

* One-time purchase
* One year of updates included
* Internal audit use permitted
* Redistribution not permitted

See `LICENSE.txt` for details.

---

## Who should use this?

* Engineering teams preparing audits
* Compliance leads
* Startup founders
* Security consultants
* Internal audit teams

If you’ve ever answered “Can you show this in GitHub?” — this tool is for you.

---