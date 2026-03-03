//go:build !wasm

package patientvisit

import (
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/time"
)

type AddMeasurementArgs struct {
	MedicalHistoryID  string  `json:"medical_history_id"`
	MeasuredByStaffID string  `json:"measured_by_staff_id"`
	MeasurementTypeID string  `json:"measurement_type_id"`
	Value             float64 `json:"value"`
	Unit              string  `json:"unit"`
	Notes             string  `json:"notes,omitempty"`
}

func (m *Module) AddMeasurement(args AddMeasurementArgs) (*ClinicalMeasurement, error) {
	if args.MedicalHistoryID == "" || args.MeasuredByStaffID == "" || args.MeasurementTypeID == "" || args.Unit == "" {
		return nil, fmt.Err("missing", "required", "arguments")
	}

	visit, err := getVisitByID(m.db, args.MedicalHistoryID)
	if err != nil {
		return nil, err
	}

	if visit.Status != StatusArrived && visit.Status != StatusTriaged && visit.Status != StatusInProgress {
		return nil, fmt.Err("visit", "not", "active")
	}

	mtype := &MeasurementType{}
	qb := m.db.Query(mtype).Where(MeasurementType_.ID).Eq(args.MeasurementTypeID)
	_, err = ReadOneMeasurementType(qb, mtype)
	if err != nil {
		return nil, fmt.Err("measurement", "type", "not", "found")
	}

	if !mtype.IsActive {
		return nil, fmt.Err("measurement", "type", "inactive")
	}

	record := &ClinicalMeasurement{
		ID:                m.uid.GetNewID(),
		MedicalHistoryID:  args.MedicalHistoryID,
		MeasuredByStaffID: args.MeasuredByStaffID,
		MeasurementTypeID: args.MeasurementTypeID,
		Value:             args.Value,
		Unit:              args.Unit,
		MeasuredAt:        time.Now(),
		Notes:             args.Notes,
	}

	if err := m.db.Create(record); err != nil {
		return nil, err
	}

	return record, nil
}

type ListMeasurementsArgs struct {
	MedicalHistoryID string `json:"medical_history_id"`
}

func (m *Module) ListMeasurements(args ListMeasurementsArgs) ([]*ClinicalMeasurement, error) {
	if args.MedicalHistoryID == "" {
		return nil, fmt.Err("missing", "medical_history_id")
	}

	qb := m.db.Query(&ClinicalMeasurement{}).
		Where(ClinicalMeasurement_.MedicalHistoryID).Eq(args.MedicalHistoryID).
		OrderBy(ClinicalMeasurement_.MeasuredAt).Asc()

	return ReadAllClinicalMeasurement(qb)
}

type ListMeasurementsByPatientArgs struct {
	PatientID         string `json:"patient_id"`
	MeasurementTypeID string `json:"measurement_type_id"`
	Limit             int    `json:"limit,omitempty"`
	Offset            int    `json:"offset,omitempty"`
}

func (m *Module) ListMeasurementsByPatient(args ListMeasurementsByPatientArgs) ([]*ClinicalMeasurement, error) {
	if args.PatientID == "" || args.MeasurementTypeID == "" {
		return nil, fmt.Err("missing", "required", "arguments")
	}

	limit := args.Limit
	if limit == 0 {
		limit = 20
	}

	// Fetch all MedicalHistory for patient
	visitsQb := m.db.Query(&MedicalHistory{}).Where(MedicalHistory_.PatientID).Eq(args.PatientID)
	visits, err := ReadAllMedicalHistory(visitsQb)
	if err != nil {
		return nil, err
	}

	if len(visits) == 0 {
		return []*ClinicalMeasurement{}, nil
	}

	visitIDs := make([]any, len(visits))
	for i, v := range visits {
		visitIDs[i] = v.ID
	}

	qb := m.db.Query(&ClinicalMeasurement{}).
		Where(ClinicalMeasurement_.MeasurementTypeID).Eq(args.MeasurementTypeID).
		Where(ClinicalMeasurement_.MedicalHistoryID).In(visitIDs).
		OrderBy(ClinicalMeasurement_.MeasuredAt).Desc().
		Limit(limit).Offset(args.Offset)

	return ReadAllClinicalMeasurement(qb)
}
