//go:build !wasm

package patientvisit_test

import (
	"testing"

	patientvisit "github.com/veltylabs/patient-visit"
)

func TestCreateVisit(t *testing.T) {
	mod, _ := setupTestModule(t)

	// OK
	args := patientvisit.CreateVisitArgs{
		PatientID:           "pat1",
		DoctorID:            "doc1",
		Reason:              "headache",
		PatientNameSnapshot: "John Doe",
		PatientRutSnapshot:  "123-4",
		DoctorNameSnapshot:  "Dr. Smith",
		AttentionAt:         1600000000,
	}

	visit, err := mod.CreateVisit(args)
	if err != nil {
		t.Fatalf("CreateVisit failed: %v", err)
	}

	if visit.Status != patientvisit.StatusCreated {
		t.Errorf("Expected status created, got %v", visit.Status)
	}

	// Missing required arg
	argsMissing := args
	argsMissing.PatientID = ""
	_, err = mod.CreateVisit(argsMissing)
	if err == nil {
		t.Errorf("Expected error for missing PatientID")
	}
}

func TestGetVisit(t *testing.T) {
	mod, _ := setupTestModule(t)
	visit, _ := mod.CreateVisit(patientvisit.CreateVisitArgs{
		PatientID:           "pat1",
		DoctorID:            "doc1",
		Reason:              "headache",
		PatientNameSnapshot: "John Doe",
		PatientRutSnapshot:  "123-4",
		DoctorNameSnapshot:  "Dr. Smith",
		AttentionAt:         1600000000,
	})

	// Found
	v, err := mod.GetVisit(patientvisit.GetVisitArgs{ID: visit.ID})
	if err != nil {
		t.Fatalf("GetVisit failed: %v", err)
	}
	if v.ID != visit.ID {
		t.Errorf("Expected visit ID %v, got %v", visit.ID, v.ID)
	}

	// Not Found
	_, err = mod.GetVisit(patientvisit.GetVisitArgs{ID: "invalid_id"})
	if err == nil {
		t.Errorf("Expected error for not found visit")
	}
}

func TestListVisitsByPatient(t *testing.T) {
	mod, _ := setupTestModule(t)

	mod.CreateVisit(patientvisit.CreateVisitArgs{
		PatientID:           "pat1",
		DoctorID:            "doc1",
		Reason:              "r1",
		PatientNameSnapshot: "p1",
		PatientRutSnapshot:  "r1",
		DoctorNameSnapshot:  "d1",
		AttentionAt:         1600000000,
	})
	mod.CreateVisit(patientvisit.CreateVisitArgs{
		PatientID:           "pat1",
		DoctorID:            "doc2",
		Reason:              "r2",
		PatientNameSnapshot: "p1",
		PatientRutSnapshot:  "r1",
		DoctorNameSnapshot:  "d2",
		AttentionAt:         1600000100,
	})

	// Found
	visits, err := mod.ListVisitsByPatient(patientvisit.ListVisitsByPatientArgs{PatientID: "pat1"})
	if err != nil {
		t.Fatalf("ListVisitsByPatient failed: %v", err)
	}
	if len(visits) != 2 {
		t.Errorf("Expected 2 visits, got %d", len(visits))
	}
	if visits[0].AttentionAt < visits[1].AttentionAt {
		t.Errorf("Expected desc order")
	}

	// Empty
	visits, err = mod.ListVisitsByPatient(patientvisit.ListVisitsByPatientArgs{PatientID: "pat_empty"})
	if err != nil {
		t.Fatalf("ListVisitsByPatient failed: %v", err)
	}
	if len(visits) != 0 {
		t.Errorf("Expected 0 visits, got %d", len(visits))
	}
}

func TestListVisitsByDoctor(t *testing.T) {
	mod, _ := setupTestModule(t)

	mod.CreateVisit(patientvisit.CreateVisitArgs{
		PatientID:           "pat1",
		DoctorID:            "doc1",
		Reason:              "r1",
		PatientNameSnapshot: "p1",
		PatientRutSnapshot:  "r1",
		DoctorNameSnapshot:  "d1",
		AttentionAt:         1600000000,
	})
	mod.CreateVisit(patientvisit.CreateVisitArgs{
		PatientID:           "pat2",
		DoctorID:            "doc1",
		Reason:              "r2",
		PatientNameSnapshot: "p2",
		PatientRutSnapshot:  "r2",
		DoctorNameSnapshot:  "d1",
		AttentionAt:         1600086390,
	})

	// All by doctor
	visits, err := mod.ListVisitsByDoctor(patientvisit.ListVisitsByDoctorArgs{DoctorID: "doc1"})
	if err != nil {
		t.Fatalf("ListVisitsByDoctor failed: %v", err)
	}
	if len(visits) != 2 {
		t.Errorf("Expected 2 visits, got %d", len(visits))
	}
	if visits[0].AttentionAt > visits[1].AttentionAt {
		t.Errorf("Expected asc order")
	}

	// Filter by date
	visitsDate, err := mod.ListVisitsByDoctor(patientvisit.ListVisitsByDoctorArgs{DoctorID: "doc1", Date: 1600000000})
	if err != nil {
		t.Fatalf("ListVisitsByDoctor by date failed: %v", err)
	}
	if len(visitsDate) != 2 {
		t.Errorf("Expected 2 visits within date range, got %d", len(visitsDate))
	}

	visitsDateOut, err := mod.ListVisitsByDoctor(patientvisit.ListVisitsByDoctorArgs{DoctorID: "doc1", Date: 1600086400})
	if err != nil {
		t.Fatalf("ListVisitsByDoctor by date failed: %v", err)
	}
	if len(visitsDateOut) != 0 {
		t.Errorf("Expected 0 visits outside date range, got %d", len(visitsDateOut))
	}
}
