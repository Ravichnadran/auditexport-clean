# AuditExport — Quick Start

Get auditor-ready GitHub evidence in **under 2 minutes**.

---

## 1️⃣ Prerequisites

* GitHub account
* Read access to the repository or organization
* GitHub Personal Access Token (PAT)

Required token scopes:

* `repo`
* `read:org`

---

## 2️⃣ Set GitHub token

```bash
export GITHUB_TOKEN=ghp_your_token_here
```

AuditExport will not run without this token.

---

## 3️⃣ Run ISO 27001 evidence collection

```bash
./auditexport run --standard iso27001
```

✔ Collects baseline GitHub evidence
✔ Skips SOC 2-only controls
✔ Safe for all repositories

---

## 4️⃣ Run SOC 2 evidence collection (extended)

```bash
./auditexport run --standard soc2
```

✔ Includes ISO 27001 evidence
✔ Adds change-management proof
✔ Adds CI/CD (GitHub Actions) evidence

---

## 5️⃣ Find your evidence

After completion:

```bash
ls evidence/
```

Final deliverable:

```
evidence/evidence.zip
```

This ZIP can be directly shared with auditors.

---

## 6️⃣ Verify what was collected

Check execution log:

```bash
cat evidence/run/execution_log.txt
```

Example:

```
run started
product standard: soc2
CI/CD evidence collected
github required reviews collected
evidence zipped
run completed
```

---

## Common questions

### ❓ Is this safe to run on production repos?

Yes. AuditExport is read-only.

### ❓ Will it fail if controls are missing?

No. Missing controls are recorded as evidence, not errors.

### ❓ Can auditors verify integrity?

Yes. SHA-256 hashes are generated for all evidence.

### ❓ Does ISO 27001 include CI/CD?

No. CI/CD evidence is included only in SOC 2 mode.

---

## That’s it

You now have **auditor-ready GitHub evidence** in minutes.

---
