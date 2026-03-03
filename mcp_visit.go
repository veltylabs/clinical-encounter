//go:build !wasm

package clinicalencounter

import (
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/orm"
	"github.com/tinywasm/time"
)

type CreateVisitArgs struct {
	PatientID               string `json:"patient_id"`
	DoctorID                string `json:"doctor_id"`
	AttentionAt             int64  `json:"attention_at"`
	Reason                  string `json:"reason"`
	PatientNameSnapshot     string `json:"patient_name_snapshot"`
	PatientRutSnapshot      string `json:"patient_rut_snapshot"`
	DoctorNameSnapshot      string `json:"doctor_name_snapshot"`
	ReservationID           string `json:"reservation_id,omitempty"`
	Diagnostic              string `json:"diagnostic,omitempty"`
	Prescription            string `json:"prescription,omitempty"`
	DoctorSpecialtySnapshot string `json:"doctor_specialty_snapshot,omitempty"`
}

func (m *Module) CreateVisit(args CreateVisitArgs) (*MedicalHistory, error) {
	if args.PatientID == "" || args.DoctorID == "" || args.Reason == "" ||
		args.PatientNameSnapshot == "" || args.PatientRutSnapshot == "" || args.DoctorNameSnapshot == "" {
		return nil, fmt.Err("missing", "required", "arguments")
	}

	if args.AttentionAt == 0 {
		return nil, fmt.Err("missing", "attention_at")
	}

	record := &MedicalHistory{
		ID:                      m.uid.GetNewID(),
		PatientID:               args.PatientID,
		DoctorID:                args.DoctorID,
		ReservationID:           args.ReservationID,
		Status:                  StatusCreated,
		AttentionAt:             args.AttentionAt,
		Reason:                  args.Reason,
		Diagnostic:              args.Diagnostic,
		Prescription:            args.Prescription,
		PatientNameSnapshot:     args.PatientNameSnapshot,
		PatientRutSnapshot:      args.PatientRutSnapshot,
		DoctorNameSnapshot:      args.DoctorNameSnapshot,
		DoctorSpecialtySnapshot: args.DoctorSpecialtySnapshot,
		UpdatedAt:               time.Now(),
	}

	if err := m.db.Create(record); err != nil {
		return nil, err
	}

	return record, nil
}

type GetVisitArgs struct {
	ID string `json:"id"`
}

func (m *Module) GetVisit(args GetVisitArgs) (*MedicalHistory, error) {
	if args.ID == "" {
		return nil, fmt.Err("missing", "id")
	}

	return getVisitByID(m.db, args.ID)
}

func getVisitByID(db *orm.DB, id string) (*MedicalHistory, error) {
	record := &MedicalHistory{}
	qb := db.Query(record).Where(MedicalHistory_.ID).Eq(id)
	_, err := ReadOneMedicalHistory(qb, record)
	if err != nil {
		return nil, fmt.Err("visit", "not", "found")
	}
	return record, nil
}

type ListVisitsByPatientArgs struct {
	PatientID string `json:"patient_id"`
	Limit     int    `json:"limit,omitempty"`
	Offset    int    `json:"offset,omitempty"`
}

func (m *Module) ListVisitsByPatient(args ListVisitsByPatientArgs) ([]*MedicalHistory, error) {
	if args.PatientID == "" {
		return nil, fmt.Err("missing", "patient_id")
	}

	limit := args.Limit
	if limit == 0 {
		limit = 20
	}

	qb := m.db.Query(&MedicalHistory{}).
		Where(MedicalHistory_.PatientID).Eq(args.PatientID).
		OrderBy(MedicalHistory_.AttentionAt).Desc().
		Limit(limit).Offset(args.Offset)

	return ReadAllMedicalHistory(qb)
}

type ListVisitsByDoctorArgs struct {
	DoctorID string `json:"doctor_id,omitempty"`
	Status   string `json:"status,omitempty"`
	Date     int64  `json:"date,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
}

func (m *Module) ListVisitsByDoctor(args ListVisitsByDoctorArgs) ([]*MedicalHistory, error) {
	limit := args.Limit
	if limit == 0 {
		limit = 20
	}

	qb := m.db.Query(&MedicalHistory{})

	if args.DoctorID != "" {
		qb = qb.Where(MedicalHistory_.DoctorID).Eq(args.DoctorID)
	}

	if args.Status != "" {
		qb = qb.Where(MedicalHistory_.Status).Eq(args.Status)
	}

	if args.Date != 0 {
		start := args.Date
		end := start + 86399
		qb = qb.Where(MedicalHistory_.AttentionAt).Gte(start).Where(MedicalHistory_.AttentionAt).Lte(end)
	}

	qb = qb.OrderBy(MedicalHistory_.AttentionAt).Asc().Limit(limit).Offset(args.Offset)

	return ReadAllMedicalHistory(qb)
}
