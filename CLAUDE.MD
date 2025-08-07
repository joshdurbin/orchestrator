# CLAUDE.md – Percona Orchestrator Project Guide

## Project Overview
- **Orchestrator**: MySQL high-availability and replication management tool, offering CLI, HTTP API, and Web UI :contentReference[oaicite:2]{index=2}.
- This is a **Percona-maintained fork** of the original project—primarily used in Percona’s Kubernetes Operators. Enhancements are published but external feature requests are not currently entertained :contentReference[oaicite:3]{index=3}.

## Workflow & Command Guide
- **Development setup**: (Add here any relevant developer setup steps, e.g., build, test commands, environment setup).
- **Core commands**:
  - `go test ./...` — run all tests
  - `go fmt ./...` — format code
  - `go vet ./...` — vet code
  - `/init` — generate initial `CLAUDE.md` skeleton (Claude Code CLI) :contentReference[oaicite:4]{index=4}.

## Code Structure & Key Packages
- `go/db/` — database connection and transaction logic :contentReference[oaicite:5]{index=5}.
- `go/config/` — parsing and validation logic for JSON config files :contentReference[oaicite:6]{index=6}.
- Include here any other major packages (e.g. HTTP API, Web UI, recovery logic).

## Style & Conventions
- Follows Go idioms: formatting via `gofmt`, error handling conventions, clear function naming.
- Configuration via JSON; ensure JSON structures are validated per schema in `go/config`.

## Testing Recommendations
- Encourage writing unit tests for new logic—especially around failover/refactoring/recovery mechanics.
- Use TDD-like workflow: “Write tests, run, then implement” :contentReference[oaicite:7]{index=7}.
- Note security-aware behavior—e.g., credential handling, auth flows (basic, multi, proxy) :contentReference[oaicite:8]{index=8}.

## Security & Access Considerations
- Be mindful of orchestrator’s handling of credentials (MySQL and HTTP basic auth)—ensure configuration is properly secured (e.g. file permissions) :contentReference[oaicite:9]{index=9}.
- For HTTP access, note the authentication modes: `basic`, `multi`, or `proxy`, and configuration patterns accordingly :contentReference[oaicite:10]{index=10}.

## Known Issues & Active Workstreams
- Current open issues include:
  - `GracefulMasterTakeover` logic flaw in semi-sync replication (opened Aug 1, 2025) :contentReference[oaicite:11]{index=11}.
  - Request IDs for API tracing (#74), adding `UseSuperReadOnly` to configurations (#68), and anonymization in UI (#49) :contentReference[oaicite:12]{index=12}.

---

### Why this matters for Claude Code
- **Project clarity**: `CLAUDE.md` gives Claude immediate context, improving task planning and accuracy :contentReference[oaicite:13]{index=13}.
- **Stylistic consistency**: Helps enforce code conventions and secure practices.
- **Faster onboarding**: Claude can assist new contributors by referencing this guide for structure, configuration, and pending workstreams.
- **Safer tooling**: Documenting CLI commands, security notes, and testing flow helps Claude suggest actions that align with team norms.

---

Let me know if you’d like this expanded with specific build/test commands, additional directory details, or tailored to a particular workflow (e.g., Kubernetes integration).
::contentReference[oaicite:14]{index=14}
