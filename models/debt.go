package models

import (
	"time"

	"gorm.io/datatypes"
)

type DebtVendor struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserId      uint      `gorm:"index;default=null" json:"user_id"`
	VendorCode  string    `gorm:"index;default=null;size=50" json:"vendor_code"`
	VendorName  string    `gorm:"default=null;size=150" json:"vendor_name"`
	Description string    `gorm:"default=null;size=200" json:"description"`
	ContactName string    `gorm:"default:null;size=50" json:"contact_name"`
	Phone       string    `gorm:"default:null;size=20" json:"phone"`
	State       int       `gorm:"index;default=1" json:"state"`
	CreatedAt   time.Time `gorm:"index;default=null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default=null" json:"updated_at"`
}

type DebtAccount struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	UserId       uint   `gorm:"index;default=null" json:"user_id"`
	VendorCode   string `gorm:"index;default=null;size=50" json:"vendor_code"`
	AccountNo    string `gorm:"index;default=null;size=100" json:"account_no"`
	AccountTitle string `gorm:"default=null;size=50" json:"account_title"` // judul akun
	DebtAmount   int    `gorm:"index;default=0" json:"debt_ammount"`
	// ===== TEMPO + TERM (misal 3 months, 6 months, 12 months) =====
	Tempo int    `gorm:"default=null" json:"tempo"`        // jumlah tempo (mis. 3, 6, 12, 24)
	Term  string `gorm:"default=null;size=30" json:"term"` // satuan tempo (daily, weekly, monthly)
	// ===== END OF TEMPO + TERM =====
	Rowstate  int       `gorm:"index;default=1" json:"rowstate"`
	CreatedAt time.Time `gorm:"index;default=null" json:"created_at"`
	UpdatedAt time.Time `gorm:"default=null" json:"updated_at"`
}

type DebtVirtualAccount struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	AccountNo string `gorm:"index;default=null;size=100" json:"account_no"`
	BankName  string `gorm:"default=null;size=50" json:"bank_name"`
	VaNumber  string `gorm:"default=null;size=50" json:"va_number"`
}

type DebtOutstanding struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	AccountNo         string         `gorm:"index;default=null;size=100" json:"account_no"`
	DebtOutstandingNo string         `gorm:"index;default=null;size=50" json:"debt_outstanding_no"`
	OutstandingAmount int64          `gorm:"default=0" json:"outstanding_amount"`
	DueDate           datatypes.Date `gorm:"default=null" json:"due_date"`
	IsPaid            bool           `gorm:"index;default=false;nullable" json:"is_paid"`
}

type DebtPayment struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	AccountNo         string         `gorm:"index;default=null;size=100" json:"account_no"`
	DebtOutstandingNo string         `gorm:"index;default=null;size=50" json:"debt_outstanding_no"`
	PaymentNo         string         `gorm:"default=null;size=50" json:"payment_no"`
	PaidAmount        int64          `gorm:"default=null" json:"paid_amount"`
	PaymentDate       datatypes.Date `gorm:"default=null" json:"payment_date"`
	CreatedAt         time.Time      `gorm:"index;default=null" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"default=null" json:"updated_at"`
}
