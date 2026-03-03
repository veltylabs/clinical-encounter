# Stage 7 — Measurement Handlers (`mcp_measurement.go`)

← [Stage 6](PLAN_STAGE_6_MEASUREMENT_TYPE.md) | Next → [Stage 8](PLAN_STAGE_8_DETAIL.md)

## Tools

| Handler | Required args | Optional args |
|---|---|---|
| `AddMeasurement` | medical_history_id, measured_by_staff_id, measurement_type_id, value, unit | notes |
| `ListMeasurements` | medical_history_id | — |
| `ListMeasurementsByPatient` | patient_id, measurement_type_id | limit (def. 20), offset (def. 0) |

## Flows

**AddMeasurement:**
1. Fetch visit by `medical_history_id`.
2. Validate `Status` ∈ {`arrived`, `triaged`, `in_progress`} — else `fmt.Err("visit","not","active")`.
   > `arrived` + `triaged` → nurse records vitals during triage.
   > `in_progress` → doctor records additional measurements.
3. Fetch `MeasurementType` by `measurement_type_id`; validate `IsActive = true` — else `fmt.Err("measurement","type","inactive")`.
4. `uid.GetNewID()` → PK; set `MeasuredAt = now()`.
5. `db.Create(&ClinicalMeasurement{...})`.
6. Return saved record.

**ListMeasurements:**
1. `db.Query(&ClinicalMeasurement{}).Where(ClinicalMeasurement_.MedicalHistoryID).Eq(id).OrderBy(ClinicalMeasurement_.MeasuredAt).Asc()`.
2. `ReadAllClinicalMeasurement(qb)` → return slice.

**ListMeasurementsByPatient:**
1. To track historical progression across multiple visits. First fetch all `MedicalHistory` (just IDs) for `patient_id`.
2. Extract the IDs from the fetched histories into `visitIDs`.
3. `db.Query(&ClinicalMeasurement{}).Where(ClinicalMeasurement_.MeasurementTypeID).Eq(measurement_type_id).Where(ClinicalMeasurement_.MedicalHistoryID).In(visitIDs).OrderBy(ClinicalMeasurement_.MeasuredAt).Desc().Limit(limit).Offset(offset)`.
4. `ReadAllClinicalMeasurement(qb)` → return slice.
