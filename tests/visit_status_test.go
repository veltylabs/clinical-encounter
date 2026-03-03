//go:build !wasm

package patientvisit_test

import (
	"testing"

	patientvisit "github.com/veltylabs/patient-visit"
)

func TestVisitStatusFSM_HappyPath(t *testing.T) {
	mod, pub := setupTestModule(t)

	// Create
	visit, _ := mod.CreateVisit(patientvisit.CreateVisitArgs{
		PatientID:           "pat1",
		DoctorID:            "doc1",
		Reason:              "checkup",
		PatientNameSnapshot: "John",
		PatientRutSnapshot:  "123",
		DoctorNameSnapshot:  "Dr. Smith",
		AttentionAt:         1600000000,
	})

	if visit.Status != patientvisit.StatusCreated {
		t.Fatalf("Expected status created, got %s", visit.Status)
	}

	// Arrived
	visit, err := mod.MarkArrived(patientvisit.MarkArrivedArgs{ID: visit.ID, PatientNameSnapshot: "John Updated"})
	if err != nil {
		t.Fatalf("MarkArrived failed: %v", err)
	}
	if visit.Status != patientvisit.StatusArrived {
		t.Errorf("Expected status arrived, got %s", visit.Status)
	}
	if visit.PatientNameSnapshot != "John Updated" {
		t.Errorf("Expected patient name to be updated")
	}
	if len(pub.events) != 1 || pub.events[0] != patientvisit.EventPatientArrived {
		t.Errorf("Expected EventPatientArrived")
	}

	// Triaged
	visit, err = mod.MarkTriaged(patientvisit.MarkTriagedArgs{ID: visit.ID})
	if err != nil {
		t.Fatalf("MarkTriaged failed: %v", err)
	}
	if visit.Status != patientvisit.StatusTriaged {
		t.Errorf("Expected status triaged, got %s", visit.Status)
	}
	if len(pub.events) != 2 || pub.events[1] != patientvisit.EventPatientTriaged {
		t.Errorf("Expected EventPatientTriaged")
	}

	// In Progress
	visit, err = mod.StartVisit(patientvisit.StartVisitArgs{ID: visit.ID})
	if err != nil {
		t.Fatalf("StartVisit failed: %v", err)
	}
	if visit.Status != patientvisit.StatusInProgress {
		t.Errorf("Expected status in_progress, got %s", visit.Status)
	}

	// Completed
	visit, err = mod.CompleteVisit(patientvisit.CompleteVisitArgs{ID: visit.ID})
	if err != nil {
		t.Fatalf("CompleteVisit failed: %v", err)
	}
	if visit.Status != patientvisit.StatusCompleted {
		t.Errorf("Expected status completed, got %s", visit.Status)
	}
	if len(pub.events) != 3 || pub.events[2] != patientvisit.EventVisitCompleted {
		t.Errorf("Expected EventVisitCompleted")
	}
}

func TestVisitStatusFSM_Cancel(t *testing.T) {
	mod, pub := setupTestModule(t)

	// From Created
	v1, _ := mod.CreateVisit(patientvisit.CreateVisitArgs{
		PatientID: "p1", DoctorID: "d1", Reason: "r1", PatientNameSnapshot: "n1", PatientRutSnapshot: "rut1", DoctorNameSnapshot: "dn1", AttentionAt: 1600000000,
	})
	v1, err := mod.CancelVisit(patientvisit.CancelVisitArgs{ID: v1.ID, Reason: "no show"})
	if err != nil {
		t.Fatalf("CancelVisit from created failed: %v", err)
	}
	if v1.Status != patientvisit.StatusCancelled {
		t.Errorf("Expected cancelled")
	}
	if len(pub.events) != 1 || pub.events[0] != patientvisit.EventVisitCancelled {
		t.Errorf("Expected EventVisitCancelled")
	}

	// Invalid transition
	_, err = mod.CancelVisit(patientvisit.CancelVisitArgs{ID: v1.ID})
	if err == nil {
		t.Errorf("Expected error cancelling an already cancelled visit")
	}

	// Complete from Arrived should fail
	v2, _ := mod.CreateVisit(patientvisit.CreateVisitArgs{
		PatientID: "p1", DoctorID: "d1", Reason: "r1", PatientNameSnapshot: "n1", PatientRutSnapshot: "rut1", DoctorNameSnapshot: "dn1", AttentionAt: 1600000000,
	})
	_, _ = mod.MarkArrived(patientvisit.MarkArrivedArgs{ID: v2.ID})
	_, err = mod.CompleteVisit(patientvisit.CompleteVisitArgs{ID: v2.ID})
	if err == nil {
		t.Errorf("Expected error completing an arrived visit")
	}
}
