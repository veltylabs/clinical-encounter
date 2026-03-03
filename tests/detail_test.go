//go:build !wasm

package patientvisit_test

import (
	"testing"

	patientvisit "github.com/veltylabs/patient-visit"
)

func TestHistoryDetail(t *testing.T) {
	mod, _ := setupTestModule(t)

	// Create visit
	visit, _ := mod.CreateVisit(patientvisit.CreateVisitArgs{
		PatientID:           "p1",
		DoctorID:            "d1",
		Reason:              "r1",
		PatientNameSnapshot: "n1",
		PatientRutSnapshot:  "rut1",
		DoctorNameSnapshot:  "dn1",
		AttentionAt:         1600000000,
	})

	// Add detail on created state (rejected)
	_, err := mod.AddHistoryDetail(patientvisit.AddHistoryDetailArgs{
		MedicalHistoryID:  visit.ID,
		CatalogItemID:     "item1",
		Quantity:          1,
		ItemNameSnapshot:  "item 1",
		ItemCodeSnapshot:  "i1",
		ItemPriceSnapshot: 10.0,
	})
	if err == nil {
		t.Errorf("Expected error adding detail on created visit")
	}

	// Move to arrived
	mod.MarkArrived(patientvisit.MarkArrivedArgs{ID: visit.ID})

	// Add detail on arrived state (rejected)
	_, err = mod.AddHistoryDetail(patientvisit.AddHistoryDetailArgs{
		MedicalHistoryID:  visit.ID,
		CatalogItemID:     "item1",
		Quantity:          1,
		ItemNameSnapshot:  "item 1",
		ItemCodeSnapshot:  "i1",
		ItemPriceSnapshot: 10.0,
	})
	if err == nil {
		t.Errorf("Expected error adding detail on arrived visit")
	}

	// Move to triaged
	mod.MarkTriaged(patientvisit.MarkTriagedArgs{ID: visit.ID})

	// Add detail on triaged state (OK)
	detail1, err := mod.AddHistoryDetail(patientvisit.AddHistoryDetailArgs{
		MedicalHistoryID:  visit.ID,
		CatalogItemID:     "item1",
		Quantity:          1,
		ItemNameSnapshot:  "item 1",
		ItemCodeSnapshot:  "i1",
		ItemPriceSnapshot: 10.0,
	})
	if err != nil {
		t.Fatalf("AddHistoryDetail failed on triaged: %v", err)
	}
	if detail1.ItemNameSnapshot != "item 1" {
		t.Errorf("Expected 'item 1', got %s", detail1.ItemNameSnapshot)
	}

	// Move to in_progress
	mod.StartVisit(patientvisit.StartVisitArgs{ID: visit.ID})

	// Add detail on in_progress state (OK)
	_, err = mod.AddHistoryDetail(patientvisit.AddHistoryDetailArgs{
		MedicalHistoryID:  visit.ID,
		CatalogItemID:     "item2",
		Quantity:          2,
		ItemNameSnapshot:  "item 2",
		ItemCodeSnapshot:  "i2",
		ItemPriceSnapshot: 20.0,
	})
	if err != nil {
		t.Fatalf("AddHistoryDetail failed on in_progress: %v", err)
	}

	// Move to completed
	mod.CompleteVisit(patientvisit.CompleteVisitArgs{ID: visit.ID})

	// Add detail on completed state (OK for billing reconciliations)
	_, err = mod.AddHistoryDetail(patientvisit.AddHistoryDetailArgs{
		MedicalHistoryID:  visit.ID,
		CatalogItemID:     "item3",
		Quantity:          3,
		ItemNameSnapshot:  "item 3",
		ItemCodeSnapshot:  "i3",
		ItemPriceSnapshot: 30.0,
	})
	if err != nil {
		t.Fatalf("AddHistoryDetail failed on completed: %v", err)
	}

	// Cancel a new visit and try to add detail
	v2, _ := mod.CreateVisit(patientvisit.CreateVisitArgs{PatientID: "p", DoctorID: "d", Reason: "r", PatientNameSnapshot: "n", PatientRutSnapshot: "rut", DoctorNameSnapshot: "dn", AttentionAt: 1600000000})
	mod.CancelVisit(patientvisit.CancelVisitArgs{ID: v2.ID, Reason: "no show"})
	_, err = mod.AddHistoryDetail(patientvisit.AddHistoryDetailArgs{
		MedicalHistoryID:  v2.ID,
		CatalogItemID:     "item1",
		Quantity:          1,
		ItemNameSnapshot:  "item 1",
		ItemCodeSnapshot:  "i1",
		ItemPriceSnapshot: 10.0,
	})
	if err == nil {
		t.Errorf("Expected error adding detail on cancelled visit")
	}

	// List history details
	details, err := mod.ListHistoryDetails(patientvisit.ListHistoryDetailsArgs{MedicalHistoryID: visit.ID})
	if err != nil {
		t.Fatalf("ListHistoryDetails failed: %v", err)
	}
	if len(details) != 3 {
		t.Errorf("Expected 3 details, got %d", len(details))
	}
}
