package cashvalidation

type UpdateSortValidation struct {
	TreasuryDetailNo string `json:"treasury_detail_no" validate:"required"`
	Sorts            int    `json:"sorts" validate:"required"`
}
