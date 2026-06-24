package models

type Treasury struct {
	InSorting int8 `gorm:"default=0" json:"in_sorting"`
}

type TreasuryMemberReq struct {
	TreasuryNo string
	Page       int
}

type InviteTreasuryMemberReq struct {
	TreasuryNo string
	Email      string
}

type AcceptInvitationReq struct {
	TreasuryNo       string `json:"treasury_no"`
	NotificationCode string `json:"notification_code"`
}

type UpdateMemberAccessReq struct {
	MemberId   int    `json:"member_id"`
	TreasuryNo string `json:"treasury_no"`
	CanEdit    int    `json:"can_edit"`
}

type RemoveMemberReq struct {
	TreasuryNo string `json:"treasury_no"`
	MemberId   int    `json:"member_id"`
}
