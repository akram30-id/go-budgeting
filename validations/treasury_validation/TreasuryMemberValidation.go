package treasuryvalidation

type ListTreasuryMembersValidation struct {
	TreasuryNo string `query:"treasury_no" json:"treasury_no" validate:"required,max=50"`
	Page       int    `query:"page" json:"page"`
}

type FindMemberValidation struct {
	Keywords string `query:"keywords" json:"keywords" validate:"required,max=50"`
}

type InviteMemberValidation struct {
	TreasuryNo string `json:"treasury_no" validate:"required,max=50"`
	Email      string `json:"email" validate:"required,max=50"`
}

type AcceptInvitationValidation struct {
	TreasuryNo       string `json:"treasury_no" validate:"required,max=50"`
	NotificationCode string `json:"notification_code" validate:"required,max=100"`
}

type MemberAccessValidation struct {
	MemberId   int    `json:"member_id" validate:"required"`
	TreasuryNo string `json:"treasury_no" validate:"required,max=50"`
	CanEdit    int    `json:"can_edit"`
}

type RemoveMemberValidation struct {
	TreasuryNo string `json:"treasury_no" validate:"required,max=50"`
	MemberId   int    `json:"member_id" validate:"required"`
}

type ListNotificationUser struct {
	Page  int `query:"page" json:"page"`
	Limit int `query:"limit" json:"limi"`
}
