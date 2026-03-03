//go:build !wasm

package clinicalencounter_test

import (
    "testing"
    "github.com/tinywasm/sqlite"
    "github.com/tinywasm/orm"
    clinicalencounter "github.com/veltylabs/clinical-encounter"
)

type mockPublisher struct{ events []string }

func (p *mockPublisher) Publish(event string, _ any) error {
    p.events = append(p.events, event)
    return nil
}

func setupTestModule(t *testing.T) (*clinicalencounter.Module, *mockPublisher) {
    t.Helper()
    db, _ := sqlite.Open(":memory:")
    for _, m := range []orm.Model{
        &clinicalencounter.MedicalHistory{}, &clinicalencounter.MeasurementType{},
        &clinicalencounter.ClinicalMeasurement{}, &clinicalencounter.HistoryDetail{},
    } {
        if err := db.CreateTable(m); err != nil {
            t.Fatalf("create table: %v", err)
        }
    }
    pub := &mockPublisher{}
    mod, err := clinicalencounter.New(db, pub)
    if err != nil {
        t.Fatalf("New: %v", err)
    }
    return mod, pub
}
