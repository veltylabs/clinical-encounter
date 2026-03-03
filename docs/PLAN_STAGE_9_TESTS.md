# Stage 9 — Tests

← [Stage 8](PLAN_STAGE_8_DETAIL.md) | Next → [Stage 10](PLAN_STAGE_10_PUBLISH.md)

## File Structure

```
tests/
├── setup_test.go               — shared setupTestModule + mockPublisher
├── visit_test.go               — CreateVisit, GetVisit, ListVisitsByPatient, ListVisitsByDoctor
├── visit_status_test.go        — FSM transitions + event assertions
├── measurement_type_test.go    — CreateMeasurementType, ListMeasurementTypes, ToggleMeasurementType
├── measurement_test.go         — AddMeasurement, ListMeasurements
└── detail_test.go              — AddHistoryDetail, ListHistoryDetails
```

## Shared setup (`tests/setup_test.go`)

```go
//go:build !wasm

package patientvisit_test

import (
    "testing"
    "github.com/tinywasm/sqlite"
    "github.com/tinywasm/orm"
    patientvisit "github.com/veltylabs/patient-visit"
)

type mockPublisher struct{ events []string }

func (p *mockPublisher) Publish(event string, _ any) error {
    p.events = append(p.events, event)
    return nil
}

func setupTestModule(t *testing.T) (*patientvisit.Module, *mockPublisher) {
    t.Helper()
    db, _ := sqlite.Open(":memory:")
    for _, m := range []orm.Model{
        &patientvisit.MedicalHistory{}, &patientvisit.MeasurementType{},
        &patientvisit.ClinicalMeasurement{}, &patientvisit.HistoryDetail{},
    } {
        if err := db.CreateTable(m); err != nil {
            t.Fatalf("create table: %v", err)
        }
    }
    pub := &mockPublisher{}
    mod, err := patientvisit.New(db, pub)
    if err != nil {
        t.Fatalf("New: %v", err)
    }
    return mod, pub
}
```

## Coverage matrix

| File | Test cases |
|---|---|
| `visit_test.go` | CreateVisit OK / missing required arg / GetVisit found / not found / ListVisitsByPatient pagination + empty / ListVisitsByDoctor filter by status + filter by date |
| `visit_status_test.go` | Full happy path (created→arrived→triaged→in_progress→completed) / EventPatientArrived published / EventPatientTriaged published / EventVisitCompleted published / CancelVisit from created, arrived, triaged, in_progress + EventVisitCancelled / invalid transition rejected |
| `measurement_type_test.go` | CreateType OK / ListTypes active-only filter / ListTypes include_inactive / Toggle activate + deactivate |
| `measurement_test.go` | AddMeasurement OK on arrived + triaged + in_progress / rejected on created + completed + cancelled / inactive type rejected / ListMeasurements empty + populated / ListMeasurementsByPatient historical tracking |
| `detail_test.go` | AddHistoryDetail OK on triaged + in_progress + completed / rejected on arrived + created + cancelled / ListHistoryDetails |
