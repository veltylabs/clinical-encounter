# FHIR Compliance — Master Plan

**Objective:** Achieve full technical FHIR compliance for the `clinical-encounter` module by implementing the three gaps identified in [`FHIR_ROADMAP.md`](FHIR_ROADMAP.md).

**Related docs:** [`ARCHITECTURE.md`](ARCHITECTURE.md) | [`SKILL.md`](SKILL.md) | [`FHIR_ROADMAP.md`](FHIR_ROADMAP.md)

---

## Development Rules

- **Agent Setup:** Run `go install github.com/tinywasm/devflow/cmd/gotest@latest` before anything.
- **Dependencies:** No external libraries. Standard library + `tinywasm/*` polyfills only (Use `tinywasm/json` instead of standard `encoding/json`).
- **Build Tags:** All server-side files MUST carry `//go:build !wasm`.
- **Database:** `model.go` is the source of truth. DO NOT edit `model_orm.go` manually; run the `ormc` code generator after changes.
- **Testing:** Split logic per rule `_back_test.go` + `//go:build !wasm`. Run via `gotest` (no args). No external assertion libraries.
- **Publishing:** Update documentation FIRST, then run `gopush 'message'`. Never `git commit/push` directly.

---

## Architecture Decision

**Hybrid: Minimal Schema Evolution + Pure Adapter Pattern**

Add domain-meaningful fields (`StartedAt`, `FinishedAt`, `Cie10Code`, `LoincCode`, `UcumUnit`) directly to the schema — these are business-valid constructs, not FHIR artifacts. A standalone FHIR adapter (`fhir_types.go` + `fhir_adapter.go`) translates internal models to FHIR R4 JSON on demand. No REST API endpoints are added; no external interfaces (`EventPublisher`, `Module.New()`) change.

---

## Remaining Tasks

The functional logic for FHIR adaptation (Stages 1 through 4) has already been pushed to the remote branch `origin/feat-fhir-compliance-3370008295037280168`.

The only remaining tasks are:

| # | Task | Status |
|---|------|--------|
| 1 | Merge the `feat-fhir-compliance...` branch into `main`. | pending |
| 2 | Pull the latest changes to the local `main` branch. | pending |
| 3 | Install/Update the newly refactored `ormc` CLI: `go install github.com/tinywasm/orm/cmd/ormc@latest` | pending |
| 4 | Run `ormc` in the root of the module to generate the new `model_orm.go` file. | pending |
| 5 | Delete the legacy `model_db.go` file (as it's been replaced by `model_orm.go`). | pending |
| 6 | Run `gotest` to ensure the project compiles and all tests pass with the new generated ORM code. | pending |

---

## Verification

```bash
# Run full test suite after completing the ORM migration
gotest

# Publish the refactored ORM integration
gopush 'chore: migrate auto-generated ORM file to model_orm.go suffix'
```
