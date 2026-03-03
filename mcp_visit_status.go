//go:build !wasm

package clinicalencounter

import (
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/time"
)

func (m *Module) applyTransition(id, action string) (*MedicalHistory, error) {
	visit, err := getVisitByID(m.db, id)
	if err != nil {
		return nil, err
	}
	next, ok := visitTransitions[visit.Status][action]
	if !ok {
		return nil, fmt.Err("invalid", "transition", visit.Status, "->", action)
	}
	visit.Status = next
	visit.UpdatedAt = time.Now()
	return visit, m.db.Update(visit)
}

type MarkArrivedArgs struct {
	ID                  string `json:"id"`
	PatientNameSnapshot string `json:"patient_name_snapshot,omitempty"`
	PatientRutSnapshot  string `json:"patient_rut_snapshot,omitempty"`
}

func (m *Module) MarkArrived(args MarkArrivedArgs) (*MedicalHistory, error) {
	if args.ID == "" {
		return nil, fmt.Err("missing", "id")
	}

	visit, err := getVisitByID(m.db, args.ID)
	if err != nil {
		return nil, err
	}

	next, ok := visitTransitions[visit.Status]["mark_arrived"]
	if !ok {
		return nil, fmt.Err("invalid", "transition", visit.Status, "->", "mark_arrived")
	}

	visit.Status = next
	visit.UpdatedAt = time.Now()

	if args.PatientNameSnapshot != "" {
		visit.PatientNameSnapshot = args.PatientNameSnapshot
	}
	if args.PatientRutSnapshot != "" {
		visit.PatientRutSnapshot = args.PatientRutSnapshot
	}

	if err := m.db.Update(visit); err != nil {
		return nil, err
	}

	m.publish(EventPatientArrived, map[string]any{
		"visit_id":              visit.ID,
		"patient_name_snapshot": visit.PatientNameSnapshot,
	})

	return visit, nil
}

type MarkTriagedArgs struct {
	ID string `json:"id"`
}

func (m *Module) MarkTriaged(args MarkTriagedArgs) (*MedicalHistory, error) {
	if args.ID == "" {
		return nil, fmt.Err("missing", "id")
	}

	visit, err := m.applyTransition(args.ID, "mark_triaged")
	if err != nil {
		return nil, err
	}

	m.publish(EventPatientTriaged, map[string]any{
		"visit_id":              visit.ID,
		"doctor_id":             visit.DoctorID,
		"patient_name_snapshot": visit.PatientNameSnapshot,
	})

	return visit, nil
}

type StartVisitArgs struct {
	ID string `json:"id"`
}

func (m *Module) StartVisit(args StartVisitArgs) (*MedicalHistory, error) {
	if args.ID == "" {
		return nil, fmt.Err("missing", "id")
	}

	return m.applyTransition(args.ID, "start_visit")
}

type CompleteVisitArgs struct {
	ID string `json:"id"`
}

func (m *Module) CompleteVisit(args CompleteVisitArgs) (*MedicalHistory, error) {
	if args.ID == "" {
		return nil, fmt.Err("missing", "id")
	}

	visit, err := m.applyTransition(args.ID, "complete")
	if err != nil {
		return nil, err
	}

	m.publish(EventVisitCompleted, map[string]any{
		"visit_id":   visit.ID,
		"patient_id": visit.PatientID,
		"doctor_id":  visit.DoctorID,
	})

	return visit, nil
}

type CancelVisitArgs struct {
	ID     string `json:"id"`
	Reason string `json:"reason,omitempty"`
}

func (m *Module) CancelVisit(args CancelVisitArgs) (*MedicalHistory, error) {
	if args.ID == "" {
		return nil, fmt.Err("missing", "id")
	}

	visit, err := m.applyTransition(args.ID, "cancel")
	if err != nil {
		return nil, err
	}

	m.publish(EventVisitCancelled, map[string]any{
		"visit_id": visit.ID,
		"reason":   args.Reason,
	})

	return visit, nil
}
