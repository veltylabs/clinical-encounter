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
