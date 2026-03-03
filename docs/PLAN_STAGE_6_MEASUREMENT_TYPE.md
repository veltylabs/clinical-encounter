# Stage 6 — MeasurementType Handlers (`mcp_measurement_type.go`)

← [Stage 5](PLAN_STAGE_5_FSM.md) | Next → [Stage 7](PLAN_STAGE_7_MEASUREMENT.md)

## Tools

| Handler | Required args | Optional args |
|---|---|---|
| `CreateMeasurementType` | name, default_unit | min_normal, max_normal |
| `ListMeasurementTypes` | — | include_inactive (bool, def. false) |
| `ToggleMeasurementType` | id, is_active | — |

## Flows

**CreateMeasurementType:**
1. `uid.GetNewID()` → PK; set `IsActive = true`.
2. `db.Create(&MeasurementType{...})`.

**ListMeasurementTypes:**
1. Base QB: `db.Query(&MeasurementType{})`.
2. If `include_inactive = false` (default): add `Where(MeasurementType_.IsActive).Eq(true)`.
3. `ReadAllMeasurementType(qb)` → return slice.

**ToggleMeasurementType:**
1. Fetch record by `id`.
2. Parse `is_active` bool from args.
3. `record.IsActive = isActive; db.Update(record)`.
