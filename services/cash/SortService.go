package cash

import (
	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpdateSortCash(req models.SortUpdate) *helpers.ReturnService {

	db := config.DB

	// AMIBIL TREASURY NO NYA DULU
	var treasuryNo string

	errGetTreasuryNo := db.Raw("SELECT treasury_no FROM treasury_detail WHERE treasury_detail_no = ? LIMIT 1", req.TreasuryDetailNo).Scan(&treasuryNo).Error

	if errGetTreasuryNo != nil {
		return &helpers.ReturnService{
			Message:  errGetTreasuryNo.Error(),
			Code:     "500SS01",
			Success:  false,
			HttpCode: 500,
		}
	}

	if treasuryNo == "" {
		return &helpers.ReturnService{
			Message:  "Treasury is not found for this cash.",
			Code:     "404F01",
			Success:  false,
			HttpCode: 404,
		}
	}

	type SimilarShorts struct {
		TreasuryDetailNo string
		Sorts            int
	}

	var similarShorts SimilarShorts

	// CEK APAKAH ADA SORTS YANG SAMA
	errGetSimilarShorts := db.Raw(`
		SELECT a.treasury_detail_no, a.sorts 
		FROM treasury_detail AS a
		WHERE a.treasury_no = (SELECT treasury_no FROM treasury_detail WHERE treasury_detail_no = ? LIMIT 1)
		AND a.sorts = ?
		LIMIT 1
	`, req.TreasuryDetailNo, req.Sorts).Scan(&similarShorts).Error

	if errGetSimilarShorts != nil {
		return &helpers.ReturnService{
			Message:  errGetSimilarShorts.Error(),
			Code:     "500SS02",
			Success:  false,
			HttpCode: 500,
		}
	}

	// jika ada sort yg sama
	if similarShorts.TreasuryDetailNo != "" {
		err := db.Transaction(func(tx *gorm.DB) error {

			// AMBIL CURRENT SORTS
			var currentSort int

			if err := tx.
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Table("treasury_detail").
				Select("sorts").
				Where("treasury_detail_no", req.TreasuryDetailNo).
				Scan(&currentSort).Error; err != nil {
				return err
			}

			// UPDATE RE-ORDER SORT
			if req.Sorts > currentSort {
				// pindah kebawah
				if err := tx.
					Table("treasury_detail").
					Where("treasury_no", treasuryNo).
					Where("sorts > ? AND sorts <= ?", currentSort, req.Sorts).
					Update("sorts", gorm.Expr("sorts - 1")).Error; err != nil {
					return err // rollback
				}

			} else if req.Sorts < currentSort {
				// pindah ke atas
				if err := tx.
					Table("treasury_detail").
					Where("treasury_no", treasuryNo).
					Where("sorts >= ? AND sorts < ?", req.Sorts, currentSort).
					Update("sorts", gorm.Expr("sorts + 1")).Error; err != nil {
					return err // rollback
				}
			}
			// end of UPDATE RE-ORDER SORT

			// UPDATE TARGET CASH
			if err := tx.
				Table("treasury_detail").
				Where("treasury_detail_no = ?", req.TreasuryDetailNo).
				Update("sorts", req.Sorts).Error; err != nil {
				return err
			}
			// end of UPDATE TARGET CASH

			return nil // commit
		})

		if err != nil {
			// transaction fail
			return &helpers.ReturnService{
				Message:  "Failed to update sorting data cash.",
				Code:     "500SS03",
				Success:  false,
				HttpCode: 500,
			}
		}
	}

	return &helpers.ReturnService{
		Message:  "Success",
		Code:     "200G00",
		Success:  true,
		HttpCode: 200,
	}

}
