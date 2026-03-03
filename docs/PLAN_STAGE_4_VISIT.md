# Stage 4 — Visit Handlers (`mcp_visit.go`)

← [Stage 3](PLAN_STAGE_3_CORE.md) | Next → [Stage 5](PLAN_STAGE_5_FSM.md)

## Tools

| Handler | Required args | Optional args |
|---|---|---|
| `CreateVisit` | patient_id, doctor_id, attention_at, reason, patient_name_snapshot, patient_rut_snapshot, doctor_name_snapshot | reservation_id, diagnostic, prescription, doctor_specialty_snapshot |
| `GetVisit` | id | — |
| `ListVisitsByPatient` | patient_id | limit (def. 20), offset (def. 0) |
| `ListVisitsByDoctor` | — | doctor_id, status, date (Unix, filters by day), limit (def. 20), offset (def. 0) |

## Flows

**CreateVisit:**
1. Parse & validate required args.
2. `uid.GetNewID()` → PK; set `Status = StatusCreated`, `UpdatedAt = now`.
3. `db.Create(&MedicalHistory{...})`.
4. Return created record.

**GetVisit:**
1. Parse `id`.
2. `m := &MedicalHistory{}; ReadOneMedicalHistory(db.Query(m).Where(MedicalHistory_.ID).Eq(id), m)`.
3. Return record or `fmt.Err("visit","not","found")`.

**ListVisitsByPatient:**
1. Parse `patient_id`, `limit` (default 20), `offset` (default 0).
2. `db.Query(&MedicalHistory{}).Where(MedicalHistory_.PatientID).Eq(patientID).OrderBy(MedicalHistory_.AttentionAt).Desc().Limit(limit).Offset(offset)`.
3. `ReadAllMedicalHistory(qb)` → return slice.

**ListVisitsByDoctor** (doctor worklist & reception pool):
1. Parse optional `doctor_id`, `status`, `date`, `limit`, `offset`.
2. Base QB: `db.Query(&MedicalHistory{})`.
3. If `doctor_id` provided: add `Where(MedicalHistory_.DoctorID).Eq(doctorID)`. (If omitted, searches globally, useful for Reception pooling `triaged` patients across all doctors).
4. If `status` provided: add `Where(MedicalHistory_.Status).Eq(status)`.
5. If `date` provided: filter `AttentionAt` ∈ [`date`, `date+86399`] (full calendar day).
6. `OrderBy(MedicalHistory_.AttentionAt).Asc()` → chronological worklist.
7. `ReadAllMedicalHistory(qb)` → return slice.
