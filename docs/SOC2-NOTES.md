# SOC 2 Notes — CI/CD & Change Management Evidence

This document explains how the attached CI/CD and GitHub evidence supports **SOC 2 Trust Services Criteria**, specifically **Change Management**, **Access Controls**, and **System Integrity**.

This document is intended for **auditors and assessors**.

---

## Scope of Evidence

The evidence in this package demonstrates controls related to:

* Change authorization
* Code review enforcement
* Build verification prior to merge
* CI/CD pipeline governance
* Tamper resistance of build processes

The system under review uses **GitHub** and **GitHub Actions** as its source control and CI/CD platform.

---

## CI/CD Overview

The organization uses **GitHub Actions** for continuous integration.

Key characteristics:

* CI/CD pipelines are defined as code (`.github/workflows/*.yml`)
* Pipelines are version-controlled
* Pipelines are triggered automatically
* Manual bypass is restricted

No external or self-hosted CI systems are used.

---

## Control: Change Authorization

### Control Objective

Ensure all code changes are reviewed and approved prior to integration.

### Evidence Provided

* Branch protection rules
* Required pull-request reviews
* Protected default branch configuration

### How to Verify

See:

* `github/protected_branches.json`
* `github/required_reviews.json`
* `github/merge_policies.json`

### Auditor Interpretation

Code cannot be merged into protected branches without satisfying approval requirements.

---

## Control: Build Verification Before Merge

### Control Objective

Ensure code changes are validated before integration.

### Evidence Provided

* GitHub Actions workflows
* Workflow trigger conditions
* Required status checks

### How to Verify

See:

* `github/cicd/workflows.json`
* `github/cicd/workflow_files/`
* `github/cicd/ci_summary.md`

### Auditor Interpretation

Build and test pipelines are automatically executed before merge completion.

---

## Control: CI/CD Pipeline Governance

### Control Objective

Ensure CI/CD processes are controlled, documented, and auditable.

### Evidence Provided

* Pipeline definitions stored in source control
* Commit history for workflow files
* Workflow execution metadata

### How to Verify

See:

* `github/cicd/workflow_files/`
* `github/cicd/workflow_runs.json`

### Auditor Interpretation

CI/CD logic changes are subject to the same change control as application code.

---

## Control: Secrets Management

### Control Objective

Ensure credentials are protected and not hardcoded.

### Evidence Provided

* Workflow references to platform-managed secrets
* Absence of plaintext credentials in workflow files

### How to Verify

Review workflow files under:

* `github/cicd/workflow_files/`

### Auditor Interpretation

Secrets are managed via GitHub’s secure secrets store and are not exposed in code.

---

## Control: Traceability & Auditability

### Control Objective

Ensure actions can be traced to users and commits.

### Evidence Provided

* Workflow run metadata
* Commit SHA references
* Actor information
* Execution timestamps

### How to Verify

See:

* `github/cicd/workflow_runs.json`
* `run/execution_log.txt`

### Auditor Interpretation

Each CI/CD execution is attributable and time-bound.

---

## What Is Explicitly Out of Scope

The following are **intentionally excluded** from evidence collection:

* Build artifacts
* Runtime logs
* Application logs
* Secret values
* Deployment execution
* Infrastructure provisioning

These are excluded to maintain **least privilege**, **data minimization**, and **audit predictability**.

---

## Evidence Integrity

All collected evidence is:

* Generated automatically
* Timestamped
* Hashed using SHA-256
* Packaged into a single ZIP archive

Integrity hashes are available in:

* `run/hashes.json`

---

## Summary for Auditors

This evidence demonstrates that:

* Code changes require authorization
* Builds run automatically prior to merge
* CI/CD processes are governed and versioned
* Secrets are securely managed
* Activities are logged and traceable

These controls align with SOC 2 requirements for **Change Management** and **System Integrity**.

---
