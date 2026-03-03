# Stage 8 — HistoryDetail Handlers (`mcp_detail.go`)

← [Stage 7](PLAN_STAGE_7_MEASUREMENT.md) | Next → [Stage 9](PLAN_STAGE_9_TESTS.md)

## Tools

| Handler | Required args | Optional args |
|---|---|---|
| `AddHistoryDetail` | medical_history_id, catalog_item_id, quantity, item_name_snapshot, item_code_snapshot, item_price_snapshot | — |
| `ListHistoryDetails` | medical_history_id | — |

## Flows

**AddHistoryDetail:**
1. Fetch visit by `medical_history_id`.
2. Validate `Status` ∈ {`triaged`, `in_progress`, `completed`} — else `fmt.Err("visit","not","active")`.
   > Items/services only after triage. Allowed in `completed` for billing reconciliations and post-consultation supply usage.
3. `uid.GetNewID()` → PK.
4. `db.Create(&HistoryDetail{...})`.
5. Return saved record.

**ListHistoryDetails:**
1. `db.Query(&HistoryDetail{}).Where(HistoryDetail_.MedicalHistoryID).Eq(id)`.
2. `ReadAllHistoryDetail(qb)` → return slice.
