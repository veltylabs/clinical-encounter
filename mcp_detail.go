//go:build !wasm

package clinicalencounter

import (
	"github.com/tinywasm/fmt"
)

type AddHistoryDetailArgs struct {
	MedicalHistoryID  string  `json:"medical_history_id"`
	CatalogItemID     string  `json:"catalog_item_id"`
	Quantity          int     `json:"quantity"`
	ItemNameSnapshot  string  `json:"item_name_snapshot"`
	ItemCodeSnapshot  string  `json:"item_code_snapshot"`
	ItemPriceSnapshot float64 `json:"item_price_snapshot"`
}

func (m *Module) AddHistoryDetail(args AddHistoryDetailArgs) (*HistoryDetail, error) {
	if args.MedicalHistoryID == "" || args.CatalogItemID == "" || args.Quantity <= 0 ||
		args.ItemNameSnapshot == "" || args.ItemCodeSnapshot == "" {
		return nil, fmt.Err("missing", "or", "invalid", "required", "arguments")
	}

	visit, err := getVisitByID(m.db, args.MedicalHistoryID)
	if err != nil {
		return nil, err
	}

	if visit.Status != StatusTriaged && visit.Status != StatusInProgress && visit.Status != StatusCompleted {
		return nil, fmt.Err("visit", "not", "active")
	}

	record := &HistoryDetail{
		ID:                m.uid.GetNewID(),
		MedicalHistoryID:  args.MedicalHistoryID,
		CatalogItemID:     args.CatalogItemID,
		Quantity:          args.Quantity,
		ItemNameSnapshot:  args.ItemNameSnapshot,
		ItemCodeSnapshot:  args.ItemCodeSnapshot,
		ItemPriceSnapshot: args.ItemPriceSnapshot,
	}

	if err := m.db.Create(record); err != nil {
		return nil, err
	}

	return record, nil
}

type ListHistoryDetailsArgs struct {
	MedicalHistoryID string `json:"medical_history_id"`
}

func (m *Module) ListHistoryDetails(args ListHistoryDetailsArgs) ([]*HistoryDetail, error) {
	if args.MedicalHistoryID == "" {
		return nil, fmt.Err("missing", "medical_history_id")
	}

	qb := m.db.Query(&HistoryDetail{}).Where(HistoryDetail_.MedicalHistoryID).Eq(args.MedicalHistoryID)
	return ReadAllHistoryDetail(qb)
}
