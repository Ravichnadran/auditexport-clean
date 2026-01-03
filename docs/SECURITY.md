# Security Overview

AuditExport is designed with a **security-first, audit-safe** architecture.

---

## Access Model

- Read-only access only
- Uses GitHub APIs with minimum required permissions
- No write, delete, or mutation operations

---

## Data Handling

- No secrets are stored
- No credentials are persisted
- No production application data is collected
- Evidence is generated locally

---

## Network Behavior

- Outbound connections only to:
  - api.github.com over HTTPS (TCP/443)
- No telemetry
- No analytics
- No third-party endpoints
- DNS and OS background traffic excluded

---

## Integrity

- Every evidence file is hashed (SHA-256)
- Any modification invalidates the evidence package
- Hashes are provided in `run/hashes.txt`

---

## Vulnerability Reporting

If you discover a security issue, report it responsibly
to the system owner or vendor contact.

Please do NOT disclose vulnerabilities publicly.
