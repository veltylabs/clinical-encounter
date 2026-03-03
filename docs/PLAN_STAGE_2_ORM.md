# Stage 2 — ORM Code (`model_db.go`)

← [Stage 1](PLAN_STAGE_1_MODELS.md) | Next → [Stage 3](PLAN_STAGE_3_CORE.md)

## Command

```bash
ormc
```

Run from the project root. It scans `model.go` and generates `model_db.go` in the same directory.

## Generated artifacts (per struct)

| Artifact | Description |
|---|---|
| `Schema() []orm.Field` | Column definitions with types and constraints |
| `Values() []any` | Field values in schema order |
| `Pointers() []any` | Field pointers for scanning rows |
| `TableName() string` | snake_case of struct name (auto) |
| `<Struct>_` | Typed column name descriptor for QB chains |
| `ReadOne<Struct>(qb, model)` | Fetch single row |
| `ReadAll<Struct>(qb)` | Fetch all rows |

## Notes

- Add build tag `//go:build !wasm` at the top of `model_db.go`.
- Do **not** declare `TableName()` manually in `model.go` — `ormc` generates it.
- Column descriptors: `MedicalHistory_`, `MeasurementType_`, `ClinicalMeasurement_`, `HistoryDetail_`.
- Use descriptors for all QB `.Where()` and `.OrderBy()` calls — never raw strings.
