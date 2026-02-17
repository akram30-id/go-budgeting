package treasury

import (
	"strconv"
	"time"

	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
	"github.com/gofiber/fiber/v2"
)

func DuplicateTreasuryService(req models.DuplicateTreasuryReq, c *fiber.Ctx) helpers.ReturnService {

	db := config.DB

	var treasuryNo string

	if err := db.
		Table("treasuries").
		Select("treasury_no").
		Where("state", 1).
		Where("treasury_no", req.TreasuryNo).
		Scan(&treasuryNo).Error; err != nil {
		return helpers.ReturnService{
			Message:  "Failed to check treasury.",
			Code:     "500E00",
			Success:  false,
			HttpCode: 500,
		}
	}

	if treasuryNo == "" {
		return helpers.ReturnService{
			Message:  "Treasury not found.",
			Code:     "404F00",
			Success:  false,
			HttpCode: 404,
		}
	}

	tx := db.Begin()

	// create new treasury

	type Treasury struct {
		ID         uint
		TreasuryNo string
		Month      string
		Year       string
		State      int
		CreatedAt  time.Time
		OwnerId    int
	}

	userID, ok := c.Locals("userId").(int)
	if !ok {
		tx.Rollback()
		return helpers.ReturnService{
			Message:  "Invalid user id.",
			Code:     "500E00",
			Success:  false,
			HttpCode: 500,
		}
	}

	dataCreateTreasury := Treasury{
		Month:     req.Month,
		Year:      req.Year,
		State:     1,
		CreatedAt: time.Now(),
		OwnerId:   userID,
	}

	if err := tx.Table("treasuries").Create(&dataCreateTreasury).Error; err != nil {
		tx.Rollback()
		return helpers.ReturnService{
			Message:  "Failed to create new treasury.",
			Code:     "500G02",
			Success:  false,
			HttpCode: 500,
		}
	}

	var generateTreasuryNo string
	generateTreasuryNo = "TRE" + helpers.StrPadLeft(strconv.Itoa(int(dataCreateTreasury.ID)), 10, "0")

	if err := tx.Table("treasuries").Where("id = ?", dataCreateTreasury.ID).Update("treasury_no", generateTreasuryNo).Error; err != nil {
		tx.Rollback()
		return helpers.ReturnService{
			Message:  "Failed to create new treasury. #2",
			Code:     "500G02",
			Success:  false,
			HttpCode: 500,
		}
	}

	// end of create new treasury

	/**
	*	create cash records
	* 	insert batch cash ke treasury baru
	 */
	for _, v := range req.TreasuryDetailNo {

		// ambil detail, expense, dan income value existing
		type ExistingTreasuryDetail struct {
			TreasuryDetailName string
			Notes              string
			IncomeValue        int64
			ExpenseValue       int64
			IsDebt             int
			Sorts              int
		}

		var existingTreasuryDetail ExistingTreasuryDetail

		if err := tx.
			Table("treasury_detail AS a").
			Select("a.treasury_detail_name, a.notes, a.income_value, a.expense_value, a.is_debt, a.sorts").
			Where("a.treasury_no", treasuryNo).
			Where("a.state", 1).
			Where("a.treasury_detail_no", v).
			Scan(&existingTreasuryDetail).Error; err != nil {

			tx.Rollback()
			return helpers.ReturnService{
				Message:  "Failed to create duplicate cash " + "(" + v + ").",
				Code:     "500G02",
				Success:  false,
				HttpCode: 500,
			}

		}

		// jika treasury detail tidak ditemukan
		if existingTreasuryDetail.TreasuryDetailName == "" {

			tx.Rollback()
			return helpers.ReturnService{
				Message:  "Cash not found " + "(" + v + ").",
				Code:     "404F00",
				Success:  false,
				HttpCode: 404,
			}

		} else { // jika ditemukan

			type TreasuryDetail struct {
				ID                 uint
				TreasuryNo         string
				TreasuryDetailNo   string
				TreasuryDetailName string
				Notes              string
				IncomeValue        int64
				ExpenseValue       int64
				IsDebt             int
				UserId             int
				Sorts              int
				CreatedAt          time.Time
			}

			treasuryDetail := TreasuryDetail{
				TreasuryNo:         generateTreasuryNo,
				Notes:              existingTreasuryDetail.Notes,
				TreasuryDetailName: existingTreasuryDetail.TreasuryDetailName,
				IncomeValue:        existingTreasuryDetail.IncomeValue,
				ExpenseValue:       existingTreasuryDetail.ExpenseValue,
				IsDebt:             existingTreasuryDetail.IsDebt,
				CreatedAt:          time.Now(),
				UserId:             userID,
				Sorts:              existingTreasuryDetail.Sorts,
			}

			if err := tx.Table("treasury_detail").Create(&treasuryDetail).Error; err != nil {

				tx.Rollback()
				return helpers.ReturnService{
					Message:  "Failed to duplicate new cash " + "(" + v + "). #1",
					Code:     "404F00",
					Success:  false,
					HttpCode: 404,
				}

			} else {
				lastId := treasuryDetail.ID

				if lastId == 0 {
					tx.Rollback()
					return helpers.ReturnService{
						Message:  "Failed to duplicate new cash: please try again " + "(" + v + "). #3",
						Code:     "500G02",
						Success:  false,
						HttpCode: 500,
					}
				}

				var generateTreasuryDetailNo string
				generateTreasuryDetailNo = "TRD" + helpers.StrPadLeft(strconv.Itoa(int(lastId)), 10, "0")

				if err := tx.Table("treasury_detail").Where("id = ?", lastId).Update("treasury_detail_no", generateTreasuryDetailNo).Error; err != nil {
					tx.Rollback()
					return helpers.ReturnService{
						Message:  "Failed to duplicate new cash: please try again " + "(" + v + "). #4",
						Code:     "500G02",
						Success:  false,
						HttpCode: 500,
					}
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return helpers.ReturnService{
			Message:  "Something's wrong. Please try again later.",
			Code:     "500E00",
			Success:  false,
			HttpCode: 500,
		}
	}

	/** end of create cash records
	 */

	return helpers.ReturnService{
		Message:  "Create new treasury with dupplicate success.",
		Code:     "200S01",
		Success:  true,
		Data:     generateTreasuryNo,
		HttpCode: 200,
	}

}
