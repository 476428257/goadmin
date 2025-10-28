package controller

import (
	"fmt"
	"strings"
	"time"

	"server/internal/middleware"
	"server/internal/model"
	"server/pkg/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminController struct {
	BaseController[model.Admin]
}

func NewAdminController() *AdminController {
	return &AdminController{
		BaseController: BaseController[model.Admin]{
			Model: model.Admin{},
		},
	}
}

func (b *AdminController) UpdateAvatar(ctx *gin.Context) {
	Model := b.Model
	id, exists := ctx.Get("userID")
	if !exists {
		b.Error(ctx, "用户ID不存在")
		return
	}
	var avatar struct {
		Avatar string `json:"avatar" binding:"required" name:"用户名"`
	}
	if err := b.HandleValidationError(ctx, &avatar); err != nil {
		return
	}
	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Model(&Model).Where("id = ?", id).Update("avatar", avatar.Avatar).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}
	b.Success(ctx, "更新成功", gin.H{
		"avatar": Model.Avatar,
	})
}

func (b *AdminController) UpdatePassword(ctx *gin.Context) {
	Model := b.Model
	id, exists := ctx.Get("userID")
	if !exists {
		b.Error(ctx, "用户ID不存在")
		return
	}
	if err := database.DB.First(&Model, id).Error; err != nil {
		b.Error(ctx, "用户不存在")
		return
	}
	var data struct {
		Old  string `json:"old" binding:"required" name:"旧密码"`
		New  string `json:"new" binding:"required" name:"新密码"`
		New1 string `json:"new1" binding:"required" name:"确认新密码"`
	}
	if err := b.HandleValidationError(ctx, &data); err != nil {
		return
	}
	if !Model.CheckPassword(data.Old) {
		b.Error(ctx, "原密码错误")
		return
	}
	if data.New != data.New1 {
		b.Error(ctx, "新密码不一致")
		return
	}
	Model.Password = data.New
	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Save(&Model).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}
	b.Success(ctx, "修改密码成功", "")
}

func (b *AdminController) List(ctx *gin.Context) {
	var admins []map[string]interface{}
	var total int64
	var pagedata struct {
		Username string `json:"username"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Pagination
	}
	if err := b.HandleValidationError(ctx, &pagedata); err != nil {
		return
	}
	prefix := CFG.MySQL.Prefix
	db := database.DB.Table(prefix + "admin a")
	if pagedata.Username != "" {
		db = db.Where("username like ?", "%"+pagedata.Username+"%")
	}
	if pagedata.Nickname != "" {
		db = db.Where("nickname like ?", "%"+pagedata.Nickname+"%")
	}
	if pagedata.Email != "" {
		db = db.Where("email like ?", "%"+pagedata.Email+"%")
	}
	if pagedata.Phone != "" {
		db = db.Where("phone like ?", "%"+pagedata.Phone+"%")
	}
	userID, _ := ctx.Get("userID")
	// 查询当前管理员的角色ID
	var roleID uint
	if err := database.DB.Model(&model.AuthUser{}).
		Where("admin_id = ?", userID).
		Select("role_id").
		Scan(&roleID).Error; err != nil {
		b.Error(ctx, "查询用户角色失败: "+err.Error())
		return
	}
	if roleID != 1 {
		db = db.Where("a.id = ?", userID)
	}
	db.Count(&total)
	offset := (pagedata.Page - 1) * pagedata.PageSize
	if err := db.
		Joins("LEFT JOIN " + prefix + "auth_user b on a.id=b.admin_id").
		Joins("LEFT JOIN " + prefix + "auth_role c on c.id=b.role_id").
		Select("a.id,username,nickname,avatar,email,a.status,created_at,updated_at,phone,b.role_id,c.title").Offset(int(offset)).Limit(int(pagedata.PageSize)).Scan(&admins).Error; err != nil {
		b.Error(ctx, err.Error())
		return
	}
	for _, admin := range admins {
		admin["created_at"] = admin["created_at"].(time.Time).Format("2006-01-02 15:04:05")
		admin["updated_at"] = admin["updated_at"].(time.Time).Format("2006-01-02 15:04:05")
	}
	b.Success(ctx, "获取列表成功", gin.H{
		"list":      &admins,
		"pageTotal": total,
	})
}

func (b *AdminController) Login(ctx *gin.Context) {
	var loginForm struct {
		Username string `json:"username" binding:"required" name:"用户名"`
		Password string `json:"password" binding:"required" name:"密码"`
	}
	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	// b.Success(ctx, "重置密码", gin.H{
	// 	"token": string(hashedPassword),
	// })
	if err := b.HandleValidationError(ctx, &loginForm); err != nil {
		return
	}
	var admin model.Admin
	// var auth_user model.AuthUser
	if err := database.DB.Where("username = ?", loginForm.Username).First(&admin).Error; err != nil {
		b.Error(ctx, "用户名不存在")
		return
	}
	if admin.Status == 0 {
		b.Error(ctx, "用户已禁用")
		return
	}
	if !admin.CheckPassword(loginForm.Password) {
		b.Error(ctx, "用户名或密码错误")
		return
	}

	// 登录更新版本号
	admin.Version++
	if admin.Version > 999 {
		admin.Version = 1
	}

	token, err := middleware.GenerateToken(admin.ID, admin.Username, admin.Version)
	if err != nil {
		b.Error(ctx, "生成令牌失败")
		return
	}
	admin.OldPassword = admin.Password
	ctx.Set("userID", admin.ID)
	ctx.Set("username", admin.Username)
	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		// 实现单点登录
		if err := tx.Model(&admin).Update("version", admin.Version).Error; err != nil {
			return fmt.Errorf("登录失败，请重新登录")
		}
		return nil
	}); err != nil {
		return
	}

	prefix := CFG.MySQL.Prefix
	rolestr := ""
	database.DB.Table(prefix+"auth_user a").
		Where("a.admin_id = ? and b.status=1", admin.ID).
		Joins("LEFT JOIN " + prefix + "auth_role b ON a.role_id = b.id").
		Select("b.rule").
		Scan(&rolestr)
	if rolestr == "" {
		b.Error(ctx, "无对应权限,登录失败")
		return
	}
	rulearr := []string{}
	if rolestr == "*" {
		rulearr = append(rulearr, rolestr)
	} else {
		rulearr = append(rulearr, strings.Split(rolestr, ",")...)
	}
	admin.Password = ""
	admin.Version = 0
	b.Success(ctx, "登录成功", gin.H{
		"token":   token,
		"admin":   &admin,
		"rulearr": rulearr,
	})
}

// Create 创建记录
func (b *AdminController) Create(ctx *gin.Context) {
	var auth_role model.AuthRole
	var auth_user model.AuthUser
	model := b.Model
	if err := b.HandleValidationError(ctx, &model); err != nil {
		return
	}
	if err := database.DB.First(&auth_role, model.RoleID).Error; err != nil {
		b.Error(ctx, "角色不存在")
		return
	}
	config, exists := ctx.Get("config")
	if !exists {
		b.Error(ctx, "配置不存在")
		return
	}
	model.Avatar = config.(map[string]string)["userhead"]
	if model.Avatar == "" {
		model.Avatar = "/image/head.png"
	}
	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(&model).Error; err != nil {
			fmt.Println(err.Error())
			return err
		}
		auth_user.RoleID = auth_role.ID
		auth_user.AdminID = model.ID
		if err := tx.Create(&auth_user).Error; err != nil {
			fmt.Println(err.Error())
			return err
		}
		return nil
	}); err != nil {
		return
	}

	b.Success(ctx, "创建成功", "")
}

func (b *AdminController) GetByID(ctx *gin.Context) {
	var Model model.Admin
	var id model.ID
	if err := b.HandleValidationError(ctx, &id); err != nil {
		return
	}
	if err := database.DB.First(&Model, id.ID).Error; err != nil {
		b.Error(ctx, "记录不存在")
		return
	}
	var auth_user model.AuthUser
	if err := database.DB.Where("admin_id = ?", Model.ID).First(&auth_user).Error; err != nil {
		b.Error(ctx, "记录不存在")
		return
	}
	Model.RoleID = auth_user.RoleID
	Model.Password = ""
	Model.Version = 0
	b.Success(ctx, "获取成功", Model)
}

func (b *AdminController) UpdateStatus(ctx *gin.Context) {
	Model := b.Model
	var data struct {
		ID     uint `json:"id" binding:"required"`
		Status uint `gorm:"column:status;" json:"status" binding:"oneof=0 1" name:"状态"`
	}
	if err := b.HandleValidationError(ctx, &data); err != nil {
		return
	}
	if err := database.DB.First(&Model, data.ID).Error; err != nil {
		b.Error(ctx, err.Error())
		return
	}
	Model.OldPassword = Model.Password
	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Model(&Model).Select("*").Updates(data).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}

	b.Success(ctx, "更新成功", "")
}

func (b *AdminController) Delete(ctx *gin.Context) {
	bmodel := b.Model
	ids := new(model.IDs)
	if err := b.HandleValidationError(ctx, &ids); err != nil {
		return
	}
	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Delete(&bmodel, "id in ?", ids.IDs).Error; err != nil {
			return err
		}
		if err := tx.Where("admin_id in ?", ids.IDs).Delete(&model.AuthUser{}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}

	b.Success(ctx, "删除成功", "")
}

func (b *AdminController) Update(ctx *gin.Context) {
	var admin model.Admin
	var updateAdmin model.UpdateAdmin
	var auth_role model.AuthRole
	var auth_user model.AuthUser
	// 绑定更新数据
	if err := b.HandleValidationError(ctx, &updateAdmin); err != nil {
		return
	}
	// 先获取原始记录
	if err := database.DB.First(&admin, updateAdmin.ID).Error; err != nil {
		b.Error(ctx, err.Error())
		return
	}
	if err := database.DB.First(&auth_role, updateAdmin.RoleID).Error; err != nil {
		b.Error(ctx, "角色不存在")
		return
	}
	admin.OldPassword = admin.Password

	if updateAdmin.Password != "" {
		admin.Password = updateAdmin.Password
	}
	admin.Nickname = updateAdmin.Nickname
	admin.RoleID = updateAdmin.RoleID
	admin.Email = updateAdmin.Email
	admin.Status = updateAdmin.Status

	if err := b.withTransaction(ctx, func(tx *gorm.DB) error {
		// 保存更新
		if err := tx.Save(&admin).Error; err != nil {
			return err
		}
		auth_user.RoleID = auth_role.ID
		if err := tx.Where("admin_id = ?", admin.ID).Updates(&auth_user).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}

	admin.Password = "" // 清除密码后再返回
	admin.Version = 0
	b.Success(ctx, "更新成功", admin)
}
