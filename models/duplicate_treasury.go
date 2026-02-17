package models

type DuplicateTreasuryReq struct {
	TreasuryNo       string   `json:"treasury_no"`
	TreasuryDetailNo []string `json:"treasury_detail_no"`
	Month            string   `json:"month"`
	Year             string   `json:"year"`
}
