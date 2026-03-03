# Stage 5 — Visit Status FSM (`mcp_visit_status.go`)

← [Stage 4](PLAN_STAGE_4_VISIT.md) | Next → [Stage 6](PLAN_STAGE_6_MEASUREMENT_TYPE.md)

## State Machine

```
created ──mark_arrived──► arrived ──mark_triaged──► triaged ──start_visit──► in_progress ──complete──► completed
   │                         │                        │
   └──cancel──► cancelled ◄──┴────────────────────────┘
```

## Tools

| Handler | Action key | Next status | Publishes |
|---|---|---|---|
| `MarkArrived` | `mark_arrived` (also accepts optional `patient_name_snapshot`, `patient_rut_snapshot` to update) | `arrived` | `EventPatientArrived` |
| `MarkTriaged` | `mark_triaged` | `triaged` | `EventPatientTriaged` → doctor |
| `StartVisit` | `start_visit` | `in_progress` | — |
| `CompleteVisit` | `complete` | `completed` | `EventVisitCompleted` |
| `CancelVisit` | `cancel` | `cancelled` | `EventVisitCancelled` |

## Shared FSM helper

```go
func (m *Module) applyTransition(id, action string) (*MedicalHistory, error) {
    visit, err := getVisitByID(m.db, id)
    if err != nil {
        return nil, err
    }
    next, ok := visitTransitions[visit.Status][action]
    if !ok {
        return nil, fmt.Err("invalid", "transition", visit.Status, "→", action)
    }
    visit.Status = next
    visit.UpdatedAt = now()
    return visit, m.db.Update(visit)
}
```

## Individual flows

**MarkArrived:**
1. Parse optional `patient_name_snapshot` and `patient_rut_snapshot` to update demographical data if it changed since creation.
2. `applyTransition(id, "mark_arrived")`. (If snapshots provided, update `visit` object before `db.Update()`).
3. `publish(EventPatientArrived, map{visit_id, patient_name_snapshot})`.

**MarkTriaged:**
1. `applyTransition(id, "mark_triaged")`.
2. `publish(EventPatientTriaged, map{visit_id, doctor_id, patient_name_snapshot})`.

**StartVisit:**
1. `applyTransition(id, "start_visit")` — no event.

**CompleteVisit:**
1. `applyTransition(id, "complete")`.
2. `publish(EventVisitCompleted, map{visit_id, patient_id, doctor_id})`.

**CancelVisit:**
1. Parse optional `reason` from args.
2. `applyTransition(id, "cancel")`.
3. `publish(EventVisitCancelled, map{visit_id, reason})`.
