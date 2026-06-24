package treasury

import (
	"strconv"
	"time"

	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/config/socket"
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
	"api-budgeting.smartcodex.cloud/services/notification"
	"github.com/gofiber/fiber/v2"
)

func ListMembersService(req *models.TreasuryMemberReq, c *fiber.Ctx) helpers.ReturnService {
	db := config.DB

	userId, ok := c.Locals("userId").(int)
	if !ok {
		return helpers.ReturnService{
			Message:  "Gagal mendapatkan informasi user.",
			Code:     "500E01",
			Success:  false,
			HttpCode: 500,
		}
	}

	type TreasuryMembers struct {
		TreasuryNo string
		MemberId   int
		State      int
		CanLook    int
		CanEdit    int
		Name       string
		Email      string
		IsAccepted int
		OwnerId    int
	}

	var treasuryMembers []TreasuryMembers
	limit := 100
	offset := 0

	if req.Page != 0 && req.Page == 1 {
		offset = 0
	}

	offset = (req.Page - 1) * limit

	if err := db.Table("treasury_members AS a").
		Joins("JOIN users AS b ON a.member_id=b.id").
		Select("a.treasury_no, a.member_id, a.state, a.can_look, a.can_edit, b.name, b.email, a.is_accepted, a.owner_id").
		Where("a.treasury_no = ?", req.TreasuryNo).
		Where("a.member_id != ?", userId).
		Where("a.state = ?", 1).Limit(limit).Offset(offset).Scan(&treasuryMembers).Error; err != nil {
		return helpers.ReturnService{
			Message:  err.Error(),
			Code:     "500E02",
			Success:  false,
			HttpCode: 500,
		}
	}

	type DataMembers struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		CanLook    int    `json:"can_look"`
		CanEdit    int    `json:"can_edit"`
		IsAccepted int    `json:"is_accepted"`
	}

	dataMembers := make([]DataMembers, 0)

	if len(treasuryMembers) > 0 {
		for _, v := range treasuryMembers {
			dataMembers = append(dataMembers, DataMembers{
				ID:         v.MemberId,
				Name:       v.Name,
				Email:      v.Email,
				CanLook:    v.CanLook,
				CanEdit:    v.CanEdit,
				IsAccepted: v.IsAccepted,
			})
		}
	}

	return helpers.ReturnService{
		Message:  "Success",
		Code:     "200S00",
		Success:  true,
		Data:     dataMembers,
		HttpCode: 200,
	}

}

func FindMemberService(keyword *string, c *fiber.Ctx) helpers.ReturnService {
	db := config.DB

	if len(*keyword) < 3 {
		return helpers.ReturnService{
			Message:  "Ketik minimal 3 karakter sabi kali...",
			Code:     "422E01",
			Success:  false,
			Data:     nil,
			HttpCode: 422,
		}
	}

	searchQuery := "%" + *keyword + "%"

	type TreasuryUsers struct {
		Name  string
		Email string
	}

	var treasuryUsers []TreasuryUsers

	if err := db.Table("users").
		Select("name, email").
		Where("name LIKE ? OR email LIKE ?", searchQuery, searchQuery).
		Where("email != ?", c.Locals("email")).
		Limit(10).Scan(&treasuryUsers).Error; err != nil {
		return helpers.ReturnService{
			Message:  "Error pas ngambil data user.",
			Code:     "500G00",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	if len(treasuryUsers) < 1 {
		return helpers.ReturnService{
			Message:  "Gak ketemu satupun.",
			Code:     "404E01",
			Success:  false,
			Data:     nil,
			HttpCode: 404,
		}
	}

	return helpers.ReturnService{
		Message:  "Success",
		Code:     "200S01",
		Success:  true,
		Data:     treasuryUsers,
		HttpCode: 200,
	}
}

func InviteMemberService(req *models.InviteTreasuryMemberReq, c *fiber.Ctx) helpers.ReturnService {

	db := config.DB

	type UserCandidate struct {
		ID    int
		Email string
	}

	var userCandidate UserCandidate

	if err := db.Table("users").Select("id, email").Where("email = ?", req.Email).First(&userCandidate).Error; err != nil {
		return helpers.ReturnService{
			Message:  "Gagal mengambil data user.",
			Code:     "500G00",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	if userCandidate.Email == "" {
		return helpers.ReturnService{
			Message:  "Usernya gak ketemu nih...",
			Code:     "404E01",
			Success:  false,
			Data:     nil,
			HttpCode: 404,
		}
	}

	var treasuryNo string

	ownerId := c.Locals("userId").(int)

	if err := db.Table("treasuries").Select("treasury_no").Where("treasury_no = ?", req.TreasuryNo).Where("state = ?", 1).Where("owner_id = ?", ownerId).Scan(&treasuryNo).Error; err != nil {
		return helpers.ReturnService{
			Message:  "Gagal saat mengambil data treasury.",
			Code:     "500G01",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	if treasuryNo == "" {
		return helpers.ReturnService{
			Message:  "Kamu Tidak Punya Akses Kesini.",
			Code:     "404E02",
			Success:  false,
			Data:     nil,
			HttpCode: 404,
		}
	}

	type TreasuryMemberCheck struct {
		MemberId   int
		isAccepted int
	}

	var treasuryMemberCheck TreasuryMemberCheck

	errIsMemberExist := db.Table("treasury_members").Select("member_id, is_accepted").Where("member_id = ?", userCandidate.ID).Where("state = ?", 1).Where("treasury_no = ?", treasuryNo).Limit(1).Scan(&treasuryMemberCheck).Error

	if errIsMemberExist != nil {
		return helpers.ReturnService{
			Message:  "Gagal mengambil data member...",
			Code:     "500G03",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	memberId := treasuryMemberCheck.MemberId
	isAccepted := treasuryMemberCheck.isAccepted

	if memberId > 0 && isAccepted == 1 {
		return helpers.ReturnService{
			Message:  "Membernya Sudah Ada Disini",
			Code:     "400E01",
			Success:  false,
			Data:     nil,
			HttpCode: 400,
		}
	}

	if memberId > 0 && isAccepted == 0 {
		return helpers.ReturnService{
			Message:  "Kamu udah pernah ngundang dia sebelumnya.",
			Code:     "400E01",
			Success:  false,
			Data:     nil,
			HttpCode: 400,
		}
	}

	invitationCode := helpers.GenerateOTP(6)

	if err := db.Table("treasury_members").Create(map[string]any{
		"treasury_no":     req.TreasuryNo,
		"member_id":       userCandidate.ID,
		"state":           1,
		"owner_id":        ownerId,
		"is_accepted":     0,
		"invitation_code": invitationCode,
		"can_look":        1,
		"can_edit":        0,
		"invited_at":      nil,
		"created_at":      time.Now(),
		"updated_at":      time.Now(),
	}).Error; err != nil {
		return helpers.ReturnService{
			Message:  "(Yah, kesalahan server). Gagal mengundang member...",
			Code:     "500G02",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	// SIMPAN NOTIFIKASI
	ownerEmail, _ := c.Locals("email").(string)
	dataNotification := models.CreateNotification{
		UserId:       userCandidate.ID,
		UserSenderId: ownerId,
		Title:        "Undangan Masuk Treasury",
		Message:      "",
		CreatedAt:    time.Now(),
		OwnerEmail:   ownerEmail,
		TreasuryNo:   treasuryNo,
	}

	createNotification := notification.CreateNotification(dataNotification)
	if !createNotification.Success {
		return createNotification
	} else {
		// Jika sukses, kirim notifikasi real-time via Hub
		targetID := strconv.Itoa(int(userCandidate.ID))
		go socket.GlobalHub.Emit(targetID, map[string]any{
			"type":              "NEW_INVITATION",
			"title":             "Undangan Masuk Treasury",
			"message":           ownerEmail + " mengundang kamu ke Treasury <b> " + treasuryNo + "</b>",
			"notification_code": createNotification.Data,
			"sender":            ownerEmail,
		})
	}

	return helpers.ReturnService{
		Message: "Yeay! berhasil mengundang " + userCandidate.Email + " kesini!",
		Code:    "200S00",
		Success: true,
		Data: map[string]any{
			"email":           userCandidate.Email,
			"invitation_code": invitationCode,
			"treasury_no":     req.TreasuryNo,
			"user_id":         userCandidate.ID,
		},
		HttpCode: 200,
	}

}

func AcceptInvitationService(req *models.AcceptInvitationReq, c *fiber.Ctx) helpers.ReturnService {
	db := config.DB

	userId, ok := c.Locals("userId").(int)
	if !ok {
		return helpers.ReturnService{
			Message:  "Gagal mendapatkan informasi user.",
			Code:     "500E01",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	// CHECK 1: Validasi invitation ada dan valid
	var invitationCheck struct {
		MemberId    int
		IsAccepted  int
		OwnerId     int
		MemberEmail string
	}

	errInvitation := db.Table("treasury_members AS a").
		Select("a.member_id, a.is_accepted, a.owner_id, b.email AS member_email").
		Joins("JOIN users AS b ON a.member_id = b.id").
		Where("member_id = ?", userId).
		Where("treasury_no = ?", req.TreasuryNo).
		Where("state = ?", 1).
		Limit(1).Scan(&invitationCheck).Error

	if errInvitation != nil {
		return helpers.ReturnService{
			Message:  "Gagal mengambil data undangan.",
			Code:     "500G01",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	if invitationCheck.MemberId == 0 {
		return helpers.ReturnService{
			Message:  "Undangan tidak ditemukan.",
			Code:     "404E01",
			Success:  false,
			Data:     nil,
			HttpCode: 404,
		}
	}

	if invitationCheck.IsAccepted == 1 {
		return helpers.ReturnService{
			Message:  "Undangan ini sudah pernah diterima sebelumnya.",
			Code:     "400E01",
			Success:  false,
			Data:     nil,
			HttpCode: 400,
		}
	}

	if invitationCheck.MemberEmail == "" {
		return helpers.ReturnService{
			Message:  "Akun tidak ditemukan.",
			Code:     "404E01",
			Success:  false,
			Data:     nil,
			HttpCode: 404,
		}
	}

	if invitationCheck.OwnerId == 0 {
		return helpers.ReturnService{
			Message:  "Owner treasury tidak ada atau invalid.",
			Code:     "400E01",
			Success:  false,
			Data:     nil,
			HttpCode: 400,
		}
	}

	// CHECK 2: Validasi notification ada dan valid
	var notificationCheck struct {
		Id    int
		State int
	}

	errNotification := db.Table("user_notifications").
		Select("id, state").
		Where("user_id = ?", userId).
		Where("notification_code = ?", req.NotificationCode).
		Limit(1).Scan(&notificationCheck).Error

	if errNotification != nil {
		return helpers.ReturnService{
			Message:  "Gagal mengambil data notifikasi.",
			Code:     "500G02",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	if notificationCheck.Id == 0 {
		return helpers.ReturnService{
			Message:  "Notifikasi tidak ditemukan.",
			Code:     "404E02",
			Success:  false,
			Data:     nil,
			HttpCode: 404,
		}
	}

	if notificationCheck.State == 0 {
		return helpers.ReturnService{
			Message:  "Notifikasi ini sudah tidak aktif.",
			Code:     "400E02",
			Success:  false,
			Data:     nil,
			HttpCode: 400,
		}
	}

	// 1. Update treasury_members: set is_accepted=1 where member_id=userId
	errUpdateMember := db.Table("treasury_members").
		Where("member_id = ?", userId).
		Where("treasury_no = ?", req.TreasuryNo).
		Where("state = ?", 1).
		Updates(map[string]any{
			"is_accepted": 1,
			"updated_at":  time.Now(),
		}).Error

	if errUpdateMember != nil {
		return helpers.ReturnService{
			Message:  "Gagal menerima undangan.",
			Code:     "500G03",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	// 2. Update user_notifications: set state=0 where user_id=userId and notification_code=req.NotificationCode
	errUpdateNotif := db.Table("user_notifications").
		Where("user_id = ?", userId).
		Where("notification_code = ?", req.NotificationCode).
		Updates(map[string]any{
			"state": 0,
		}).Error

	if errUpdateNotif != nil {
		return helpers.ReturnService{
			Message:  "Gagal update notifikasi.",
			Code:     "500G04",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	notificationMessage := `
				<div class="col-12 list-notif rounded-2 mb-1" style="height: 128px; width: 380px;">
					<div class="row align-items-center">
						<div class="col-2 mt-3">
							<i class="bi bi-person-fill" style="font-size: 3rem;"></i>
						</div>
						<div class="col-10">
							` + invitationCheck.MemberEmail + " menerima undangan masuk ke treasury " + req.TreasuryNo + `
						</div>
					</div>
				</div>`

	createNotification := notification.CreateNotification(models.CreateNotification{
		UserId:       invitationCheck.OwnerId,
		UserSenderId: userId,
		OwnerEmail:   "",
		TreasuryNo:   req.TreasuryNo,
		Title:        "Undangan Treasury Diterima",
		Message:      notificationMessage,
	})

	if !createNotification.Success {
		return createNotification
	} else {
		// Jika sukses, kirim notifikasi real-time via Hub
		targetID := strconv.Itoa(invitationCheck.OwnerId)
		go socket.GlobalHub.Emit(targetID, map[string]any{
			"type":              "INVITATION_ACCEPTED",
			"title":             "Undangan Masuk Treasury",
			"message":           createNotification.Message,
			"notification_code": createNotification.Data,
			"sender":            invitationCheck.MemberId,
		})
	}

	return helpers.ReturnService{
		Message:  "Berhasil menerima undangan!",
		Code:     "200S00",
		Success:  true,
		Data:     nil,
		HttpCode: 200,
	}
}

func UpdateMemberAccess(req *models.UpdateMemberAccessReq, c *fiber.Ctx) helpers.ReturnService {

	db := config.DB

	userId, ok := c.Locals("userId").(int)
	if !ok {
		return helpers.ReturnService{
			Message:  "Gagal mendapatkan informasi user.",
			Code:     "500E01",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	type TreasuryMember struct {
		MemberId   int    `json:"member_id"`
		OwnerId    int    `json:"owner_id"`
		TreasuryNo string `json:"treasury_no"`
		CanEdit    int    `json:"can_edit"`
		CanLook    int    `json:"can_look"`
		IsAccepted int    `json:"is_accepted"`
		State      int    `json:"state"`
	}

	var treasuryMember TreasuryMember

	errCheckMemberExist := db.Table("treasury_members").
		Select("member_id, owner_id, treasury_no, can_edit, can_look, is_accepted, state").
		Where("member_id = ?", req.MemberId).
		Where("treasury_no = ?", req.TreasuryNo).
		Where("state = ?", 1).
		Where("is_accepted = ?", 1).
		Limit(1).
		Scan(&treasuryMember).
		Error

	if errCheckMemberExist != nil {
		return helpers.ReturnService{
			Message:  "Gagal mengambil informasi member.",
			Code:     "500E01",
			Success:  false,
			HttpCode: 500,
		}
	}

	if treasuryMember.MemberId == 0 {
		return helpers.ReturnService{
			Message:  "Member tidak ditemukan.",
			Code:     "404E01",
			Success:  false,
			HttpCode: 404,
		}
	}

	if treasuryMember.CanEdit == req.CanEdit {
		return helpers.ReturnService{
			Message: "Berhasil mengubah akses member.",
			Code:    "200S01",
			Success: true,
			Data:    treasuryMember,
		}
	}

	if userId != treasuryMember.OwnerId {
		return helpers.ReturnService{
			Message: "Tindakan ini dilarang.",
			Code:    "401E01",
			Success: false,
		}
	}

	errUpdateMemberAccess := db.Table("treasury_members").
		Where("member_id = ?", treasuryMember.MemberId).
		Where("treasury_no = ?", treasuryMember.TreasuryNo).
		Update("can_edit", req.CanEdit).
		Error

	if errUpdateMemberAccess != nil {
		return helpers.ReturnService{
			Message:  "Gagal mengubah akses member.",
			Code:     "500E02",
			Success:  false,
			HttpCode: 500,
		}
	}

	treasuryMember.CanEdit = req.CanEdit

	return helpers.ReturnService{
		Message: "Berhasil mengubah akses member.",
		Code:    "200S02",
		Success: true,
		Data:    treasuryMember,
	}

}

func RemoveMemberTreasury(c *fiber.Ctx, req models.RemoveMemberReq) helpers.ReturnService {

	db := config.DB

	userId, ok := c.Locals("userId").(int)
	if !ok {
		return helpers.ReturnService{
			Message:  "Gagal mendapatkan informasi user.",
			Code:     "500E01",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	var ownerId int

	// check is treasury exist and active
	getTreasury := db.Table("treasuries AS a").Select("a.owner_id").Where("treasury_no = ?", req.TreasuryNo).Where("state = ?", 1).Scan(&ownerId)

	if getTreasury.Error != nil {
		return helpers.ReturnService{
			Message:  "Terdapat kesalahan (UNEXPECTED_001).",
			Code:     "500E01",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	if ownerId == 0 {
		return helpers.ReturnService{
			Message:  "Treasury tidak ditemukan.",
			Code:     "404E01",
			Success:  false,
			Data:     nil,
			HttpCode: 404,
		}
	}

	// check is user have ownership
	if userId != ownerId {
		return helpers.ReturnService{
			Message:  "Akses Tidak Valid.",
			Code:     "401E01",
			Success:  false,
			Data:     nil,
			HttpCode: 401,
		}
	}

	// softdelete member
	deleteMember := db.Table("treasury_members AS a").
		Where("a.treasury_no = ?", req.TreasuryNo).
		Where("a.member_id = ?", req.MemberId).
		Update("state", 0)

	if deleteMember.Error != nil {
		return helpers.ReturnService{
			Message:  "Terdapat kesalahan (UNEXPECTED_002).",
			Code:     "500E02",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	return helpers.ReturnService{
		Message: "Berhasil Menghapus Member.",
		Code:    "200S01",
		Success: true,
	}

}
