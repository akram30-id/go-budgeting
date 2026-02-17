package treasuryvalidation

type DuplicateTreasuryValidation struct {
	TreasuryNo       string   `json:"treasury_no" validate:"required,max=50"`
	TreasuryDetailNo []string `json:"treasury_detail_no" validate:"required"`
	Month            string   `json:"month" validate:"required,max=10"`
	Year             string   `json:"year" validate:"required,max=10"`
}
