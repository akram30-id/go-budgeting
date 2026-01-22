package services

import (
	"crypto/md5"
	"encoding/hex"

	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
)

func RegisterUser(register models.ReqRegisterUser) helpers.ReturnService {

	db := config.DB

	var existingEmail string

	err := db.Raw("SELECT email FROM users WHERE email = ? LIMIT 1", register.Email).Scan(&existingEmail).Error

	if err != nil {
		return helpers.ReturnService{
			Message:  "Gagal mengambil data dari server.",
			Code:     "500G00",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	if existingEmail != "" {
		return helpers.ReturnService{
			Message:  "User sudah ada.",
			Code:     "400F01",
			Success:  false,
			Data:     nil,
			HttpCode: 400,
		}
	}

	// validate role
	role := roleCheck(register.RoleId)

	if !role.Success {
		return helpers.ReturnService{
			Message:  role.Message,
			Code:     "400F01",
			Success:  false,
			Data:     nil,
			HttpCode: 400,
		}
	}

	tx := db.Begin()

	if tx.Error != nil {
		return helpers.ReturnService{
			Message:  "Gagal membuat transaksi ke database.",
			Code:     "500G00",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	if register.Password != register.ConfirmPassword {
		return helpers.ReturnService{
			Message:  "Konfirmasi Password tidak cocok.",
			Code:     "422F01",
			Success:  false,
			Data:     nil,
			HttpCode: 422,
		}
	}

	hash := md5.New()
	hash.Write([]byte(register.Password))

	register.Password = hex.EncodeToString(hash.Sum(nil))

	if err := tx.Exec("INSERT INTO users (email, password, name, role_id) VALUES (?,?,?,?)", register.Email, register.Password, register.Name, register.RoleId).Error; err != nil {
		tx.Rollback()
		return helpers.ReturnService{
			Message:  "Registrasi gagal.",
			Code:     "500F01",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	tx.Commit()

	return helpers.ReturnService{
		Message: "Registrasi berhasil.",
		Code:    "200F01",
		Success: true,
		Data: map[string]any{
			"email": register.Email,
			"role":  role.Data.Role,
		},
		HttpCode: 200,
	}
}

// check is role valid
func roleCheck(roleId int) models.ReturnRole {

	db := config.DB

	var role models.Roles

	err := db.Raw("SELECT id, role FROM roles AS a WHERE a.id = ? AND state = 1", roleId).Scan(&role).Error

	if err != nil {
		return models.ReturnRole{
			Success: false,
			Message: "Gagal menapatkan user role.",
		}
	}

	if role.Role == "" {
		return models.ReturnRole{
			Success: false,
			Message: "Role tidak valid.",
		}
	}

	return models.ReturnRole{
		Success: true,
		Message: "Success",
		Data:    role,
	}

}
