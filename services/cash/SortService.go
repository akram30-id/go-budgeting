package cash

import (
	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
	"gorm.io/gorm"
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

	// IS IN SORTING
	errInSorting, isInSorting := isInSorting(treasuryNo)
	if errInSorting != nil {
		return &helpers.ReturnService{
			Message:  errInSorting.Error(),
			Code:     "503E00",
			Success:  false,
			HttpCode: 503,
		}
	}

	if isInSorting == true {
		return &helpers.ReturnService{
			Message:  "Lagi ada sorting yang diproses, tunggu sebentar lalu coba lagi ya...",
			Code:     "400E00",
			Success:  false,
			HttpCode: 400,
		}
	}

	// LOCK IN_SORTING TREASURY NYA BAIR GA BENTROK
	err, locking := lockSorting(treasuryNo)
	if err != nil {
		return &helpers.ReturnService{
			Message:  locking,
			Code:     "500E02",
			Success:  false,
			HttpCode: 500,
		}
	}

	err, msgNormalizeOrder := normalizeSortOrder(treasuryNo)
	if err != nil {

		// RELEASE LOCK
		err, releaseLock := releaseLock(treasuryNo)
		if err != nil {
			return &helpers.ReturnService{
				Message:  releaseLock,
				Code:     "500E02",
				Success:  false,
				HttpCode: 422,
			}
		}

		return &helpers.ReturnService{
			Message:  msgNormalizeOrder,
			Code:     "500E03",
			Success:  false,
			HttpCode: 500,
		}
	}

	// type SimilarShorts struct {
	// 	TreasuryDetailNo string
	// 	Sorts            int
	// }

	// var similarShorts SimilarShorts

	// // CEK APAKAH ADA SORTS YANG SAMA
	// errGetSimilarShorts := db.Raw(`
	// 	SELECT a.treasury_detail_no, a.sorts
	// 	FROM treasury_detail AS a
	// 	WHERE a.treasury_no = (SELECT treasury_no FROM treasury_detail WHERE treasury_detail_no = ? LIMIT 1)
	// 	AND a.sorts = ?
	// 	LIMIT 1
	// `, req.TreasuryDetailNo, req.Sorts).Scan(&similarShorts).Error

	// if errGetSimilarShorts != nil {
	// 	return &helpers.ReturnService{
	// 		Message:  errGetSimilarShorts.Error(),
	// 		Code:     "500SS02",
	// 		Success:  false,
	// 		HttpCode: 500,
	// 	}
	// }

	// // jika ada sort yg sama
	// if similarShorts.TreasuryDetailNo != "" {

	// }

	trxErr := db.Transaction(func(tx *gorm.DB) error {

		// AMBIL CURRENT SORTS
		var currentSort int

		if err := tx.
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

	if trxErr != nil {
		// transaction fail

		err, releaseLock := releaseLock(treasuryNo)
		if err != nil {
			return &helpers.ReturnService{
				Message:  releaseLock,
				Code:     "500E02",
				Success:  false,
				HttpCode: 422,
			}
		}

		return &helpers.ReturnService{
			Message:  "Failed to update sorting data cash.",
			Code:     "500SS03",
			Success:  false,
			HttpCode: 500,
		}
	}

	err, releaseLock := releaseLock(treasuryNo)
	if err != nil {
		return &helpers.ReturnService{
			Message:  releaseLock,
			Code:     "500E02",
			Success:  false,
			HttpCode: 422,
		}
	}

	return &helpers.ReturnService{
		Message:  "Success",
		Code:     "200G00",
		Success:  true,
		HttpCode: 200,
	}

}

func isInSorting(treasuryNo string) (error, bool) {

	db := config.DB

	var inSorting int8

	row := db.Table("treasuries").
		Select(`COALESCE(in_sorting,0)`).
		Where("treasury_no = ?", treasuryNo).
		Where("state = ?", 1).
		Row()

	if err := row.Scan(&inSorting); err != nil {
		return err, false
	}

	if inSorting == 1 {
		return nil, true
	}

	return nil, false
}

func lockSorting(treasuryNo string) (error, string) {
	db := config.DB

	err, inSorting := isInSorting(treasuryNo)
	if err != nil {
		return err, "Gagal mengambil data sorting."
	}

	if inSorting == true {
		return nil, "Sedang ada proses sorting."
	}

	if err := db.Table("treasuries").Where("treasury_no = ?", treasuryNo).Where("state = ?", 1).Update("in_sorting", 1).Error; err != nil {
		return err, "[LOCK_FAIL] Yah, sortingnya gagal. Coba lagi nanti ya..."
	}

	return nil, "Sukses melakukan sorting."
}

func releaseLock(treasuryNo string) (error, string) {
	db := config.DB

	if err := db.Table("treasuries").Where("treasury_no = ?", treasuryNo).Where("state = ?", 1).Update("in_sorting", 0).Error; err != nil {
		return err, "[RELEASE_FAIL] Yah, sortingnya gagal. Coba lagi nanti ya..."
	}

	return nil, "Sukses melakukan sorting."
}

func normalizeSortOrder(treasuryNo string) (error, string) {
	db := config.DB

	type Cashes struct {
		TreasuryDetailNo string
	}

	var cash []Cashes

	if err := db.Table("treasury_detail").
		Select("treasury_detail_no").
		Where("treasury_no = ?", treasuryNo).
		Where("state = ?", 1).Order("sorts ASC").Scan(&cash).Error; err != nil {
		return err, "Gagal mendapatkan data cash."
	}

	if len(cash) < 1 {
		return nil, "Tidak ada data cash ditemukan."
	}

	txErr := db.Transaction(func(tx *gorm.DB) error {
		for k, v := range cash {
			if err := tx.Table("treasury_detail").
				Where("treasury_detail_no = ?", v.TreasuryDetailNo).
				Update("sorts", k+1).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return txErr, "Gagal normalisasi list cash."
	}

	return nil, "Normalisasi berhasil."

}
