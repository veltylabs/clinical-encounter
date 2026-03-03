//go:build !wasm

package patientvisit_test

import (
	"testing"

	patientvisit "github.com/veltylabs/patient-visit"
)

func TestMeasurement(t *testing.T) {
	mod, _ := setupTestModule(t)

	// Setup data
	mtype, _ := mod.CreateMeasurementType(patientvisit.CreateMeasurementTypeArgs{Name: "BP", DefaultUnit: "mmHg"})
	visit, _ := mod.CreateVisit(patientvisit.CreateVisitArgs{
		PatientID:           "p1",
		DoctorID:            "d1",
		Reason:              "r1",
		PatientNameSnapshot: "n1",
		PatientRutSnapshot:  "rut1",
		DoctorNameSnapshot:  "dn1",
		AttentionAt:         1600000000,
	})

	// Add measurement on created state (should be rejected)
	_, err := mod.AddMeasurement(patientvisit.AddMeasurementArgs{
		MedicalHistoryID:  visit.ID,
		MeasuredByStaffID: "nurse1",
		MeasurementTypeID: mtype.ID,
		Value:             120.0,
		Unit:              "mmHg",
	})
	if err == nil {
		t.Errorf("Expected error adding measurement to created visit")
	}

	// Move to arrived -> triaged
	mod.MarkArrived(patientvisit.MarkArrivedArgs{ID: visit.ID})
	mod.MarkTriaged(patientvisit.MarkTriagedArgs{ID: visit.ID})

	// Add measurement on triaged state (should be OK)
	m1, err := mod.AddMeasurement(patientvisit.AddMeasurementArgs{
		MedicalHistoryID:  visit.ID,
		MeasuredByStaffID: "nurse1",
		MeasurementTypeID: mtype.ID,
		Value:             120.0,
		Unit:              "mmHg",
	})
	if err != nil {
		t.Fatalf("AddMeasurement failed: %v", err)
	}
	if m1.Value != 120.0 {
		t.Errorf("Expected 120.0, got %f", m1.Value)
	}

	// Move to in_progress
	mod.StartVisit(patientvisit.StartVisitArgs{ID: visit.ID})

	// Add measurement on in_progress state (should be OK)
	_, err = mod.AddMeasurement(patientvisit.AddMeasurementArgs{
		MedicalHistoryID:  visit.ID,
		MeasuredByStaffID: "doc1",
		MeasurementTypeID: mtype.ID,
		Value:             130.0,
		Unit:              "mmHg",
	})
	if err != nil {
		t.Fatalf("AddMeasurement failed: %v", err)
	}

	// Move to completed
	mod.CompleteVisit(patientvisit.CompleteVisitArgs{ID: visit.ID})

	// Add measurement on completed state (should be rejected)
	_, err = mod.AddMeasurement(patientvisit.AddMeasurementArgs{
		MedicalHistoryID:  visit.ID,
		MeasuredByStaffID: "nurse1",
		MeasurementTypeID: mtype.ID,
		Value:             140.0,
		Unit:              "mmHg",
	})
	if err == nil {
		t.Errorf("Expected error adding measurement to completed visit")
	}

	// Inactive type test
	inactiveType, _ := mod.CreateMeasurementType(patientvisit.CreateMeasurementTypeArgs{Name: "InactiveType", DefaultUnit: "unit"})
	mod.ToggleMeasurementType(patientvisit.ToggleMeasurementTypeArgs{ID: inactiveType.ID, IsActive: false})

	visit2, _ := mod.CreateVisit(patientvisit.CreateVisitArgs{PatientID: "p1", DoctorID: "d1", Reason: "r1", PatientNameSnapshot: "n1", PatientRutSnapshot: "rut1", DoctorNameSnapshot: "dn1", AttentionAt: 1600000000})
	mod.MarkArrived(patientvisit.MarkArrivedArgs{ID: visit2.ID})

	_, err = mod.AddMeasurement(patientvisit.AddMeasurementArgs{
		MedicalHistoryID:  visit2.ID,
		MeasuredByStaffID: "nurse1",
		MeasurementTypeID: inactiveType.ID,
		Value:             1.0,
		Unit:              "unit",
	})
	if err == nil {
		t.Errorf("Expected error for inactive measurement type")
	}

	// List Measurements
	measurements, err := mod.ListMeasurements(patientvisit.ListMeasurementsArgs{MedicalHistoryID: visit.ID})
	if err != nil {
		t.Fatalf("ListMeasurements failed: %v", err)
	}
	if len(measurements) != 2 {
		t.Errorf("Expected 2 measurements, got %d", len(measurements))
	}

	// List Measurements By Patient
	patientMeasurements, err := mod.ListMeasurementsByPatient(patientvisit.ListMeasurementsByPatientArgs{
		PatientID:         "p1",
		MeasurementTypeID: mtype.ID,
	})
	if err != nil {
		t.Fatalf("ListMeasurementsByPatient failed: %v", err)
	}
	if len(patientMeasurements) != 2 {
		t.Errorf("Expected 2 patient measurements, got %d", len(patientMeasurements))
	}
}
