package model

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Admin struct {
	Commonmodel
	Username    string `gorm:"column:username;size:32;uniqueIndex;not null" json:"username" binding:"required,min=4,max=32,alphanum" name:"用户名"`
	Password    string `gorm:"column:password;size:128;not null" json:"password,omitempty" binding:"required,min=6,max=32" name:"密码"`
	OldPassword string `gorm:"-" json:"-"`
	Nickname    string `gorm:"column:nickname;size:32" json:"nickname" binding:"required,min=2,max=32" name:"昵称"`
	Email       string `gorm:"column:email;size:64;" json:"email" binding:"required,email" name:"邮箱"`
	Phone       string `gorm:"column:phone;size:11;" json:"phone" binding:"required" name:"手机号"`
	Avatar      string `gorm:"column:avatar;" json:"avatar"`
	RoleID      uint   `gorm:"-" json:"role_id,omitempty" binding:"required" name:"权限角色"`
	Version     uint   `gorm:"column:version;autoIncrement" json:"version,omitempty"`
	Status      int    `gorm:"column:status;default:1" json:"status" binding:"required,oneof=0 1" name:"状态"`
}

type UpdateAdmin struct {
	ID       uint   `json:"id" binding:"required" name:"用户id"`
	RoleID   uint   `gorm:"-" json:"role_id,omitempty" binding:"required" name:"权限角色"`
	Password string `json:"password" binding:"omitempty,min=6,max=32" name:"密码"`
	Nickname string `json:"nickname" binding:"required,min=2,max=32" name:"昵称"`
	Email    string `json:"email" binding:"required,email" name:"邮箱"`
	Status   int    `json:"status" binding:"oneof=0 1" name:"状态"`
}

// ValidatePassword 自定义密码验证规则
func (a *Admin) ValidatePassword() error {
	// 密码复杂度验证

	// hasNumber   = regexp.MustCompile(`[0-9]`).MatchString(a.Password)
	// hasLower    = regexp.MustCompile(`[a-z]`).MatchString(a.Password)
	// hasUpper    = regexp.MustCompile(`[A-Z]`).MatchString(a.Password)
	// hasSpecial  = regexp.MustCompile(`[!@#$%^&*]`).MatchString(a.Password)
	validLength := len(a.Password) >= 6 && len(a.Password) <= 32
	// hasUsername = strings.Contains(strings.ToLower(a.Password), strings.ToLower(a.Username))

	if !validLength {
		return errors.New("密码长度必须在6-32位之间")
	}
	// if !hasNumber {
	// 	return errors.New("密码必须包含数字")
	// }
	// if !hasLower {
	// 	return errors.New("密码必须包含小写字母")
	// }
	// if !hasUpper {
	// 	return errors.New("密码必须包含大写字母")
	// }
	// if !hasSpecial {
	// 	return errors.New("密码必须包含特殊字符!@#$%^&*")
	// }
	// if hasUsername {
	// 	return errors.New("密码不能包含用户名")
	// }

	return nil
}

func (a *Admin) BeforeSave(*gorm.DB) error {
	if a.Password != a.OldPassword {
		if err := a.ValidatePassword(); err != nil {
			return err
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		a.Password = string(hashedPassword)
	}
	return nil
}

// 自定义时间格式序列化
// func (a *Admin) MarshalJSON() ([]byte, error) {
// 	type Alias Admin
// 	return json.Marshal(&struct {
// 		*Alias
// 		CreatedAt string `json:"created_at"`
// 		UpdatedAt string `json:"updated_at"`
// 	}{
// 		Alias:     (*Alias)(a),
// 		CreatedAt: a.CreatedAt.Format("2006-01-02 15:04:05"),
// 		UpdatedAt: a.UpdatedAt.Format("2006-01-02 15:04:05"),
// 	})
// }

// HashPassword 加密密码
func (a *Admin) HashPassword() error {
	if a.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		a.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword 检查密码是否正确
func (a *Admin) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}
