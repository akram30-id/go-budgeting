package auth

import (
	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func ChangePasswordService(req models.ChangePasswordRequest, c *fiber.Ctx) helpers.ReturnService {

	db := config.DB

	email, ok := c.Locals("email").(string)
	if !ok {
		return helpers.ReturnService{
			Message:  "Invalid user.",
			Code:     "500E00",
			Success:  false,
			HttpCode: 500,
		}
	}

	// cek apakah old passwordnya bener atau engga

	type UserResult struct {
		ID       int
		Email    string
		Password string
		Name     string
		RoleId   int
	}

	var userResult UserResult

	if err := db.Table("users").
		Select("id", "email", "password", "name", "role_id").
		Where("email = ?", email).
		Scan(&userResult).Error; err != nil {
		return helpers.ReturnService{
			Message:  "Gagal saat mengambil data kamu.",
			Code:     "500G02",
			Success:  false,
			HttpCode: 500,
		}
	}

	if userResult.ID == 0 {
		return helpers.ReturnService{
			Message:  "User kamu gak ketemu nih.",
			Code:     "401F01",
			Success:  false,
			HttpCode: 404,
		}
	}

	// VERIFY HASHED PASSWORD
	// jika salah
	if err := bcrypt.CompareHashAndPassword([]byte(userResult.Password), []byte(req.OldPassword)); err != nil {
		return helpers.ReturnService{
			Message:  "Password lama nya gak bener tau. Cek lagi deh...",
			Code:     "400F01",
			Success:  false,
			HttpCode: 400,
		}
	}

	// cek apakah password baru dan konfirmasinya sama
	if req.NewPassword != req.ConfirmPassword {
		return helpers.ReturnService{
			Message:  "Dih, password baru sama konfirmasi password nya gak sama.",
			Code:     "422E01",
			Success:  false,
			HttpCode: 422,
		}
	}

	// jika benar
	// ubah password

	// hashed new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return helpers.ReturnService{
			Message:  "Gagal saat enkripsi password baru kamu. Coba lagi aja ya...",
			Code:     "500G01",
			Success:  false,
			HttpCode: 500,
		}
	}

	if err := db.Table("users").Where("email = ?", email).Update("password", string(hashedPassword)).Error; err != nil {
		return helpers.ReturnService{
			Message:  "Update password gagal nih. Coba lagi kali ya...",
			Code:     "500G03",
			Success:  false,
			HttpCode: 500,
		}
	}

	// return success
	return helpers.ReturnService{
		Message:  "OK",
		Code:     "200S0",
		Success:  true,
		Data:     "Yeay! password berhasil di-update.",
		HttpCode: 200,
	}
}
