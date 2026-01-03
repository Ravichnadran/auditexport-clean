# AuditExport

AuditExport is a **CLI-based audit evidence generator** for engineering teams preparing for  
**ISO/IEC 27001** and **SOC 2** audits.

It collects **read-only, verifiable, auditor-ready technical evidence**  
from GitHub and packages it into a structured, hash-verified archive.

AuditExport focuses on **evidence generation**, not certification, dashboards, or monitoring.

---

## What AuditExport Is

- A **local CLI tool** (runs on your machine or CI)
- Generates **technical audit evidence**, not audit opinions
- Uses **read-only GitHub API access**
- Produces **deterministic, timestamped output**
- Designed for **auditors, founders, security teams, and consultants**

---

## What AuditExport Is NOT

- ❌ Not a SOC 2 or ISO certification service
- ❌ Not a compliance guarantee
- ❌ Not a SaaS or hosted platform
- ❌ No agents, no background monitoring
- ❌ No write access to your repositories
- ❌ No telemetry, tracking, or analytics

---

## Supported Standards & Editions

### ISO/IEC 27001 — Community Edition (Free)

Collects governance and access-control evidence from GitHub:

- Organization and repository inventory
- Branch configuration
- Commit and pull request history
- Code ownership
- Access controls and permissions

Suitable for **ISO 27001 readiness and internal reviews**.

---

### SOC 2 — Professional Edition

Includes everything in ISO 27001 Community Edition, plus:

- Enforced change management evidence
- Pull request enforcement
- Independent review validation
- CI/CD execution proof (GitHub Actions)
- PR → CI → merge correlation
- Failed CI negative evidence
- Reviewer independence checks
- SOC 2 control mapping (CC6, CC7, CC8)
- Auditor-readable SOC 2 assertions

---

## Installation & Usage

```bash
# Make binary executable
chmod +x auditexport

# Export read-only GitHub token
export GITHUB_TOKEN=<your_read_only_token>

# Run ISO/IEC 27001 (Community Edition)
auditexport run --standard iso27001

# Run SOC 2 (Professional Edition)
auditexport run --standard soc2

# Run SOC 2 with an audit time window
auditexport run --standard soc2 \
  --from-date 2025-10-01 \
  --to-date   2025-12-31
```

---

## GitHub Access Requirements

Required token scopes:

- `repo:read`
- `read:org`
- `read:actions` (required only for SOC 2 CI/CD evidence)

AuditExport never modifies repositories or organization settings.

---

## Evidence Output Structure

```
evidence/
├── run/
│   ├── run_metadata.json
│   ├── execution_log.txt
│   └── hashes.txt
│
├── github/
│   ├── organization.json
│   ├── repositories.json
│   ├── branches.json
│   ├── commits.json
│   ├── pull_requests.json
│   ├── contributors.json
│   ├── access_controls.json
│   ├── protected_branches.json
│   └── workflows/
│       ├── workflows.json
│       └── workflow_runs.json
│
├── summaries/
│   ├── executive_summary.txt
│   ├── technical_summary.txt
│   ├── auditor_notes.txt
│   ├── soc2_change_management_assertions.txt
│   ├── soc2_extended_assertions.txt
│   └── soc2_control_mapping.json
│
└── evidence.zip
```

---

## Evidence Integrity

- Every file is hashed using **SHA-256**
- Hashes are stored in `run/hashes.txt`
- Any modification invalidates the evidence package
- Evidence is generated locally and deterministically

---

## Handling Missing or Skipped Evidence

If a control is not enabled, not applicable, or not accessible due to permissions:

- The condition is explicitly recorded
- No placeholder or fabricated evidence is created
- This is a valid and transparent audit state

---

## For Auditors

- All evidence is collected **read-only**
- No secrets or production data are exported
- Raw evidence and interpreted summaries are clearly separated
- Hashes can be used to verify integrity

See **AUDITOR_GUIDE.md** for detailed review instructions.

---

## License

AuditExport is distributed under a **commercial license**.  
See the `LICENSE` file for details.

---

## Disclaimer

AuditExport provides **technical audit evidence only**.  
Final audit opinions, conclusions, and certifications  
remain the responsibility of the auditor.
